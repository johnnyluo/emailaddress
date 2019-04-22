package emailaddress

import (
	"bytes"
	"fmt"
)

const (
	// MaxLocalPart is the maximum length of the local part
	MaxLocalPart = 64
	// MaxDomainLength the total length of domain should be less than 255 characters
	MaxDomainLength        = 255
	specialLocalCharacters = ` ",:;<>@[\]`
	validLocalPartChars    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!#$%&'*+-/=?^_`{|}~;.()"
	validDomainChars       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-"
)

var (
	// ErrEmptyEmail when the given email is actually empty
	ErrEmptyEmail = fmt.Errorf("empty string is not valid email address")
)

// email represent an email address
type email struct {
	localPartContainsSpecialChar bool
	localPart                    []rune
	domain                       []rune
}

func (e email) String() string {
	return fmt.Sprintf("%s@%s", string(e.localPart), string(e.domain))
}

// Validate the given email address
func Validate(emailAddress string) (bool, error) {
	if len(emailAddress) == 0 {
		return false, ErrEmptyEmail
	}

	localPart := make([]rune, 0, len(emailAddress))
	domain := make([]rune, 0, len(emailAddress))
	seeAt := false
	localPartContainsSpecialChar := false
	seeQuotation := false
	var previousChar rune
	for _, item := range emailAddress {
		switch item {
		case '"':
			if previousChar == '\\' {
				previousChar = '"'
				localPart = append(localPart, item)
				continue
			}

			if seeQuotation {
				seeQuotation = false
			} else {
				seeQuotation = true
			}
			if seeAt {
				return false, fmt.Errorf("%c is invalid character in domain part", item)
			}
			localPart = append(localPart, item)
			previousChar = '"'
		case '@':
			if seeQuotation {
				// '@' is valid inside quoted string
				previousChar = '@'
				localPart = append(localPart, item)
				continue
			}
			if seeAt {
				// means there are multiple '@' in the email address
				return false, fmt.Errorf("an email address can't have multiple '@' characters")
			}
			seeAt = true
			previousChar = '@'
		default:
			if !seeAt {
				isValidLocalPartChar := bytes.ContainsRune([]byte(validLocalPartChars), item)
				isSpecialLocalPartChar := bytes.ContainsRune([]byte(specialLocalCharacters), item)
				if !isValidLocalPartChar && !isSpecialLocalPartChar {
					return false, fmt.Errorf("%c is not a valid character in the local part of an email address", item)
				}

				if isSpecialLocalPartChar && !seeQuotation {
					return false, fmt.Errorf("%c only valid inside quoted string", item)
				}
				localPart = append(localPart, item)
			} else {
				domain = append(domain, item)
			}
			previousChar = item
		}
	}
	if !seeAt {
		return false, fmt.Errorf("%s is not valid email address, the format of email addresses is local-part@domain", emailAddress)
	}
	if len(localPart) == 0 {
		return false, fmt.Errorf("email address can't start with '@'")
	}
	if len(localPart) > MaxLocalPart {
		return false, fmt.Errorf("the length of local part should be less than %d", MaxLocalPart)
	}
	if len(domain) > MaxDomainLength {
		return false, fmt.Errorf("%s is longer than %d", string(domain), MaxDomainLength)
	}
	if len(domain) == 0 {
		return false, fmt.Errorf("domain part can't be empty")
	}
	e := email{
		localPartContainsSpecialChar: localPartContainsSpecialChar,
		localPart:                    localPart,
		domain:                       domain,
	}
	localPartResult, err := e.validateLocalPart()
	if nil != err {
		return localPartResult, err
	}
	return true, nil
}

// validateLocalPart of the email address
func (e email) validateLocalPart() (bool, error) {
	if len(e.localPart) == 0 {
		return false, fmt.Errorf("")
	}
	fmt.Println(string(e.localPart))
	quotation := false
	lastChar := e.localPart[len(e.localPart)-1]
	if e.localPart[0] == '"' && lastChar == '"' {
		quotation = true
	}
	if e.localPartContainsSpecialChar && !quotation {
		return false, fmt.Errorf("%s are only allowed inside a quoted string", specialLocalCharacters)
	}

	lp := e.localPart
	if quotation {
		lp = e.localPart[1 : len(e.localPart)-1]
	}
	var previousChar rune
	firstBackSlash := true
	for _, item := range lp {
		switch item {
		case '\\':
			if firstBackSlash {
				firstBackSlash = false
				previousChar = '\\'
				continue
			}
			firstBackSlash = true
		case '"':
			if previousChar != '\\' {
				return false, fmt.Errorf("\" need to be preceded by a backslash")
			}
		case '.':
			// consective dot only valid inside quotation
			if previousChar == '.' && !quotation {
				return false, fmt.Errorf("consective dot only valid inside quotation")
			}
		}
		previousChar = item
	}
	return true, nil
}
