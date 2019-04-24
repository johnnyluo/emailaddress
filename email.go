package emailaddress

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/pkg/errors"
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
	lp     *localPart
	domain []rune
}

// localPart represent the localpart of an email address
type localPart struct {
	comment           string
	localPartEmail    string
	tags              []string
	commentAtBegining bool
}

func (e email) String() string {
	return fmt.Sprintf("%s@%s", e.lp, string(e.domain))
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

	lp := make([]rune, 0, len(input))
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
			lp = append(lp, item)
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
		return nil, errors.Wrap(err, "fail to parse localPart of the email address")
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

	localp := make([]rune, 0, localPartLength)
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
