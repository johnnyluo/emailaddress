package emailaddress // import "github.com/johnnyluo/emailaddress"

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
)

var (
	// ErrEmptyEmail when the given email is actually empty
	ErrEmptyEmail = fmt.Errorf("empty string is not valid email address")
	byteEscape    = byte('\\')
	// ErrInvalidLocalPart indicate the local part of email is invalid
	ErrInvalidLocalPart = fmt.Errorf("invalid local part")
)

// email represent an email address
type email struct {
	lp     *localPart
	domain []byte
}

func (e email) String() string {
	return fmt.Sprintf("%s@%s", e.lp, string(e.domain))
}

// localPart represent the localpart of an email address
type localPart struct {
	comment           string
	localPartEmail    string
	tags              []string
	commentAtBegining bool
}

// String convert the local part back
func (lp localPart) String() string {
	b := strings.Builder{}
	b.Reset()
	if len(lp.comment) > 0 && lp.commentAtBegining {
		b.WriteString("(" + lp.comment + ")")
	}
	b.WriteString(lp.localPartEmail)
	for _, t := range lp.tags {
		b.WriteString("+" + t)
	}
	if len(lp.comment) > 0 && !lp.commentAtBegining {
		b.WriteString("(" + lp.comment + ")")
	}
	return b.String()
}

// Validate the given email address
func Validate(emailAddress string) (bool, error) {

	if len(emailAddress) == 0 {
		return false, ErrEmptyEmail
	}
	_, err := parseEmailAddress(emailAddress)
	if nil != err {
		return false, err
	}
	return true, nil
}

// parseEmailAddress
func parseEmailAddress(input string) (*email, error) {
	if len(input) == 0 {
		return nil, ErrEmptyEmail
	}

	lp := make([]byte, 0, len(input))
	domain := make([]byte, 0, len(input))
	seeAt := false
	inQuotation := false
	var previousChar byte
	for i := 0; i < len(input); i++ {
		c := input[i]
		switch c {
		case '"':
			if previousChar != byteEscape {
				inQuotation = !inQuotation
			}
		case '@':
			if !inQuotation && previousChar != byteEscape {
				if seeAt {
					// means there are multiple '@' in the email address
					return nil, fmt.Errorf("an email address can't have multiple '@' characters")
				}
				seeAt = true
				continue
			}
		}
		previousChar = c
		if seeAt {
			domain = append(domain, c)
		} else {
			lp = append(lp, c)
		}
	}

	if !seeAt {
		return nil, fmt.Errorf("%s is not valid email address, the format of email addresses is local-part@domain", input)
	}
	if len(lp) == 0 {
		return nil, fmt.Errorf("email address can't start with '@'")
	}
	if len(lp) > MaxLocalPart {
		return nil, fmt.Errorf("the length of local part should be less than %d", MaxLocalPart)
	}
	if len(domain) > MaxDomainLength {
		return nil, fmt.Errorf("%s is longer than %d", string(domain), MaxDomainLength)
	}
	if len(domain) == 0 {
		return nil, fmt.Errorf("domain part can't be empty")
	}
	lpp, err := parseLocalPart(string(lp))
	if nil != err {
		return nil, fmt.Errorf("fail to parse localPart of the email address")
	}
	if !IsDomainName(string(domain)) {
		return nil, fmt.Errorf("%s is not a valid domain", string(domain))
	}

	return &email{
		lp:     lpp,
		domain: domain,
	}, nil
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
				localPartEmail: lp,
			}, nil
		}
		return nil, fmt.Errorf("%s is invalid in the local part of an email address", lp)
	}

	localp := make([]byte, 0, localPartLength)
	inQuotation := false
	var previousChar byte
	escape := 0
	seeTag := false
	commentStart := -1
	commentEnd := -1
	for idx := 0; idx < len(lp); idx++ {
		c := lp[idx]
		switch c {
		case '"':
			if previousChar != byteEscape {
				inQuotation = !inQuotation
			}
		case '+':
			if previousChar != byteEscape {
				seeTag = true
			}
		case '.':
			if idx == 0 || idx == (localPartLength-1) {
				return nil, fmt.Errorf("%c can't be the start or end of local part", c)
			}
			if previousChar == '.' && !inQuotation {
				return nil, fmt.Errorf("consective dot is only valid in quotation")
			}
		case byteEscape:
			escape++
		case ',', ':', ';', '<', '>', '@', '[', ']', ' ':
			if !inQuotation && previousChar != byteEscape {
				return nil, fmt.Errorf("%c is only valid in quoted string or escaped", c)
			}
		case '(':
			if !inQuotation && previousChar != byteEscape {
				commentStart = idx
			}
		case ')':
			if !inQuotation && previousChar != byteEscape {
				commentEnd = idx
			}
		default:
			if previousChar == byteEscape && !inQuotation {
				return nil, fmt.Errorf("\\ is only valid in quoted string or escaped")
			}
		}

		if escape > 0 && escape%2 == 0 {
			previousChar = 0
		} else {
			previousChar = c
		}
		if !seeTag && (commentStart < 0 || (commentEnd > 0 && idx > commentEnd)) {
			// legitimate localpart
			localp = append(localp, c)
		}
	}
	if inQuotation {
		return nil, fmt.Errorf("\" is only valid escaped with baskslash")
	}
	if commentStart > -1 && commentEnd == -1 {
		return nil, fmt.Errorf("( is only valid within quoted string or escaped")
	}
	if commentStart == -1 && commentEnd > -1 {
		return nil, fmt.Errorf(") is only valid within quoted string or escaped")
	}
	lpResult := &localPart{
		localPartEmail: string(localp),
	}
	if seeTag {
		lpResult.tags = strings.Split(lp, "+")[1:]
	}

	if commentStart > -1 && commentEnd > -1 {
		lpResult.comment = string(lp[commentStart+1 : commentEnd])
		if commentStart == 0 {
			lpResult.commentAtBegining = true
		}
	}
	return lpResult, nil
}

// Equals will parse the given email addresses , and then compare it.
// if first or second are not legitimate email address, this function will return false
func Equals(first string, second string) bool {
	eFirst, err := parseEmailAddress(first)
	if nil != err {
		return false
	}
	eSec, err := parseEmailAddress(second)
	if nil != err {
		return false
	}
	if !bytes.EqualFold(eFirst.domain, eSec.domain) {
		return false
	}
	if !strings.EqualFold(eFirst.lp.localPartEmail, eSec.lp.localPartEmail) {
		return false
	}
	return true
}
