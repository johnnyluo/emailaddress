package emailaddress

import (
	"bytes"
	"fmt"
	"strings"
)

const (
	// MaxLocalPart is the maximum length of the local part
	MaxLocalPart = 64
	// MaxDomainLength the total length of domain should be less than 255 characters
	MaxDomainLength        = 255
	specialLocalCharacters = ` ",:;<>@[\]`
	validLocalPartChars    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!#$%&'*+-/=?^_`{|}~;"
	validDomainChars       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-"
)

var (
	// ErrEmptyEmail when the given email is actually empty
	ErrEmptyEmail = fmt.Errorf("empty string is not valid email address")
)

// email represent an email address
type email struct {
	localPart []rune
	domain    []rune
}

// localPart represent the localpart of an email address
type localPart struct {
	comment   string
	localPart string
	tags      []string
}

func (e email) String() string {
	return fmt.Sprintf("%s@%s", string(e.localPart), string(e.domain))
}

// Validate the given email address
func Validate(emailAddress string) (bool, error) {
	return true, nil
	// if len(emailAddress) == 0 {
	// 	return false, ErrEmptyEmail
	// }

	// localPart := make([]rune, 0, len(emailAddress))
	// domain := make([]rune, 0, len(emailAddress))
	// seeAt := false
	// seeQuotation := false
	// var previousChar rune
	// for _, item := range emailAddress {
	// 	switch item {
	// 	case '"':
	// 		if previousChar == '\\' {
	// 			previousChar = '"'
	// 			localPart = append(localPart, item)
	// 			continue
	// 		}

	// 		if seeQuotation {
	// 			seeQuotation = false
	// 		} else {
	// 			seeQuotation = true
	// 		}
	// 		if seeAt {
	// 			return false, fmt.Errorf("%c is invalid character in domain part", item)
	// 		}
	// 		localPart = append(localPart, item)
	// 		previousChar = '"'
	// 	case '@':
	// 		if seeQuotation {
	// 			// '@' is valid inside quoted string
	// 			previousChar = '@'
	// 			localPart = append(localPart, item)
	// 			continue
	// 		}
	// 		if seeAt {
	// 			// means there are multiple '@' in the email address
	// 			return false, fmt.Errorf("an email address can't have multiple '@' characters")
	// 		}
	// 		seeAt = true
	// 		previousChar = '@'
	// 	default:
	// 		if !seeAt {
	// 			localPart = append(localPart, item)
	// 		} else {
	// 			domain = append(domain, item)
	// 		}
	// 		previousChar = item
	// 	}
	// }
	// if !seeAt {
	// 	return false, fmt.Errorf("%s is not valid email address, the format of email addresses is local-part@domain", emailAddress)
	// }
	// if len(localPart) == 0 {
	// 	return false, fmt.Errorf("email address can't start with '@'")
	// }
	// if len(localPart) > MaxLocalPart {
	// 	return false, fmt.Errorf("the length of local part should be less than %d", MaxLocalPart)
	// }
	// if len(domain) > MaxDomainLength {
	// 	return false, fmt.Errorf("%s is longer than %d", string(domain), MaxDomainLength)
	// }
	// if len(domain) == 0 {
	// 	return false, fmt.Errorf("domain part can't be empty")
	// }
	// e := email{
	// 	localPart: localPart,
	// 	domain:    domain,
	// }
	// localPartResult, err := e.validateLocalPart()
	// if nil != err {
	// 	return localPartResult, err
	// }
	// return true, nil
}

// parseEmailAddress
func parseEmailAddress(input string) (*email, error) {
	if len(input) == 0 {
		return nil, ErrEmptyEmail
	}

	localPart := make([]rune, 0, len(input))
	domain := make([]rune, 0, len(input))

	seeAt := false
	inQuotation := false
	var previousChar rune

	for _, item := range input {
		switch item {
		case '"':
			if previousChar != '\\' {
				inQuotation = !inQuotation
			}
		case '@':
			if !inQuotation && previousChar != '\\' {
				if seeAt {
					// means there are multiple '@' in the email address
					return nil, fmt.Errorf("an email address can't have multiple '@' characters")
				}
				seeAt = true
				continue
			}
		}
		previousChar = item
		if seeAt {
			domain = append(domain, item)
		} else {
			localPart = append(localPart, item)
		}
	}

	if !seeAt {
		return nil, fmt.Errorf("%s is not valid email address, the format of email addresses is local-part@domain", input)
	}
	if len(localPart) == 0 {
		return nil, fmt.Errorf("email address can't start with '@'")
	}
	if len(localPart) > MaxLocalPart {
		return nil, fmt.Errorf("the length of local part should be less than %d", MaxLocalPart)
	}
	if len(domain) > MaxDomainLength {
		return nil, fmt.Errorf("%s is longer than %d", string(domain), MaxDomainLength)
	}
	if len(domain) == 0 {
		return nil, fmt.Errorf("domain part can't be empty")
	}
	fmt.Printf("local:%s\n", string(localPart))
	fmt.Printf("domain:%s\n", string(domain))
	return nil, nil
}

// parseLocalPart of email address
func parseLocalPart(lp string) (*localPart, error) {
	localPartLength := len(lp)
	if localPartLength == 0 {
		return nil, fmt.Errorf("empty local part")
	}
	// special case , local part only has one character
	if localPartLength == 1 {
		if bytes.Contains([]byte(validLocalPartChars), []byte(lp)) {
			return &localPart{
				localPart: lp,
			}, nil
		}
		return nil, fmt.Errorf("%s is invalid in the local part of an email address", lp)
	}

	localp := make([]rune, localPartLength)
	inQuotation := false
	var previousChar rune
	escape := 0
	seeTag := false
	commentStart := -1
	commentEnd := -1

	for idx, item := range lp {
		switch item {
		case '"':
			if previousChar != '\\' {
				inQuotation = !inQuotation
			}
		case '+':
			if previousChar != '\\' {
				seeTag = true
			}
		case '.':
			if idx == 0 || idx == (localPartLength-1) {
				return nil, fmt.Errorf("%c can't be the start or end of local part", item)
			}
			if previousChar == '.' && !inQuotation {
				return nil, fmt.Errorf("consective dot is only valid in quotation")
			}
		case '\\':
			escape++
		case ',', ':', ';', '<', '>', '@', '[', ']', ' ':
			if !inQuotation && previousChar != '\\' {
				return nil, fmt.Errorf("%c is only valid in quoted string or escaped", item)
			}
		case '(':
			if !inQuotation && previousChar != '\\' {
				commentStart = idx
			}
		case ')':
			if !inQuotation && previousChar != '\\' {
				commentEnd = idx
			}
		default:
			if previousChar == '\\' && !inQuotation {
				return nil, fmt.Errorf("\\ is only valid in quoted string or escaped")
			}
		}

		if escape > 0 && escape%2 == 0 {
			previousChar = -1
		} else {
			previousChar = item
		}
		if !seeTag && (commentStart < 0 || (commentEnd > 0 && idx > commentEnd)) {
			// legitimate localpart
			localp = append(localp, item)
		}
	}
	if inQuotation {
		return nil, fmt.Errorf("\" is only valid escaped with baskslash")
	}
	if commentStart > 0 && commentEnd == -1 {
		return nil, fmt.Errorf("( is only valid within quoted string or escaped")
	}
	if commentStart == -1 && commentEnd > 0 {
		return nil, fmt.Errorf(") is only valid within quoted string or escaped")
	}
	lpResult := &localPart{
		localPart: string(localp),
	}
	if seeTag {
		lpResult.tags = strings.Split(lp, "+")[1:]
	}

	if commentStart > -1 && commentEnd > -1 {
		lpResult.comment = string(lp[commentStart:commentEnd])
	}
	return lpResult, nil

}

// // validateLocalPart of the email address
// func (e email) validateLocalPart() (bool, error) {
// 	if len(e.localPart) == 0 {
// 		return false, fmt.Errorf("")
// 	}
// 	quotation := false
// 	lastChar := e.localPart[len(e.localPart)-1]
// 	if e.localPart[0] == '"' && lastChar == '"' {
// 		quotation = true
// 	}

// 	lp := e.localPart
// 	if quotation {
// 		lp = e.localPart[1 : len(e.localPart)-1]
// 	}
// 	var previousChar rune
// 	firstBackSlash := true
// 	for _, item := range lp {
// 		switch item {

// 		case '\\':
// 			if firstBackSlash {
// 				firstBackSlash = false
// 				previousChar = '\\'
// 				continue
// 			}
// 			firstBackSlash = true
// 		case '"':
// 			if previousChar != '\\' {
// 				return false, fmt.Errorf("\" need to be preceded by a backslash")
// 			}
// 		case '.':
// 			// consective dot only valid inside quotation
// 			if previousChar == '.' && !quotation {
// 				return false, fmt.Errorf("consective dot only valid inside quotation")
// 			}
// 		}
// 		previousChar = item
// 	}
// 	return true, nil
// }
