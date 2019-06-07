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

// tag represent tag in email local part
type tag struct {
	emailTags string
	start     int
	end       int
}

// String stringer implementation
func (t tag) String() string {
	totalLen := len(t.emailTags)
	if t.end == 0 || len(t.emailTags) == 0 || t.start > totalLen || t.end > totalLen || t.end < t.start {
		return ""
	}
	return t.emailTags[t.start:t.end]
}

// email represent an email address
type email struct {
	lp     *localPart
	domain string
}

func (e email) String() string {
	return fmt.Sprintf("%s@%s", e.lp, string(e.domain))
}

// localPart represent the localpart of an email address
type localPart struct {
	comment           string
	localPartEmail    string
	tags              []tag
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
		b.WriteString("+" + t.String())
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

	atLoc := -1
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
				atLoc = i
				continue
			}
		}
		previousChar = c
	}

	if !seeAt {
		return nil, fmt.Errorf("%s is not valid email address, the format of email addresses is local-part@domain", input)
	}

	lenDomain := len(input) - atLoc - 1
	if atLoc == 0 {
		return nil, fmt.Errorf("email address can't start with '@'")
	}
	if atLoc > MaxLocalPart {
		return nil, fmt.Errorf("the length of local part should be less than %d", MaxLocalPart)
	}
	if lenDomain > MaxDomainLength {
		return nil, fmt.Errorf("%s is longer than %d", string(input[atLoc+1:]), MaxDomainLength)
	}
	if lenDomain == 0 {
		return nil, fmt.Errorf("domain part can't be empty")
	}
	lpp, err := parseLocalPart(string(input[:atLoc]))
	if nil != err {
		return nil, fmt.Errorf("fail to parse localPart of the email address")
	}
	if !IsDomainName(string(input[atLoc+1:])) {
		return nil, fmt.Errorf("%s is not a valid domain", string(input[atLoc+1:]))
	}

	return &email{
		lp:     lpp,
		domain: input[atLoc+1:],
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

	inQuotation := false
	var previousChar byte
	escape := 0
	seeTag := false
	commentStart := -1
	commentEnd := -1
	start := 0
	end := localPartLength

	for idx := 0; idx < len(lp); idx++ {
		c := lp[idx]
		switch c {
		case '"':
			if previousChar != byteEscape {
				inQuotation = !inQuotation
			}
		case '+':
			if previousChar != byteEscape {
				if !seeTag {
					end = idx
				}
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

	if commentStart > commentEnd {
		return nil, fmt.Errorf("invalid email address")
	}
	if commentStart == 0 {
		start = commentEnd + 1
	}
	if commentStart > 0 {
		end = commentStart
	}

	lpResult := &localPart{
		localPartEmail: lp[start:end],
	}
	if seeTag {
		lpResult.tags = getTags(lp[end:])
	}

	if commentStart > -1 && commentEnd > -1 {
		lpResult.comment = string(lp[commentStart+1 : commentEnd])
		if commentStart == 0 {
			lpResult.commentAtBegining = true
		}
	}
	return lpResult, nil
}

func getTags(t string) []tag {
	totalLen := len(t)
	if totalLen == 0 || t == "+" {
		return nil
	}
	var tags []tag
	currentTag := tag{
		emailTags: t,
		start:     0,
	}
	for idx := 0; idx < totalLen; idx++ {
		c := t[idx]
		switch c {
		case '+':
			if idx == 0 {
				// it start with +, we skip it
				currentTag.start = 1
				continue
			}
			currentTag.end = idx
			tags = append(tags, currentTag)
			currentTag = tag{
				emailTags: t,
				start:     idx + 1,
			}
		}
	}
	currentTag.end = totalLen
	if currentTag.start < currentTag.end {
		// if it is end with tag, then
		tags = append(tags, currentTag)
	}
	return tags
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
	if !strings.EqualFold(eFirst.domain, eSec.domain) {
		return false
	}
	if !strings.EqualFold(eFirst.lp.localPartEmail, eSec.lp.localPartEmail) {
		return false
	}
	return true
}
