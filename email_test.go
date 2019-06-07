package emailaddress

import (
	"fmt"
	"reflect"
	"regexp"
	"testing"
)

func TestValidate(t *testing.T) {
	cases := []struct {
		name           string
		input          string
		expectedResult bool
		err            error
	}{
		{
			name:           "empty-email",
			input:          "",
			expectedResult: false,
			err:            ErrEmptyEmail,
		},
		{
			name:           "single-quote email",
			input:          `"@test.net`,
			expectedResult: false,
			err:            fmt.Errorf(`"@test.net is not valid email address, the format of email addresses is local-part@domain`),
		},
		{
			name:           "double-quote email",
			input:          `"we\"d"@test.net`,
			expectedResult: true,
			err:            nil,
		},
		{
			name:           "consective dot email",
			input:          `we..johnny@test.net`,
			expectedResult: false,
			err:            fmt.Errorf("fail to parse localPart of the email address"),
		},
		{
			name:           "consective dot email",
			input:          `"we..johnny"@test.net`,
			expectedResult: true,
			err:            nil,
		},
		{
			name:           "email with comment",
			input:          `john.smith(comment)@example.com`,
			expectedResult: true,
			err:            nil,
		},
		{
			name:           "email with comment1",
			input:          `(comment)john.smith@example.com`,
			expectedResult: true,
			err:            nil,
		},
		{
			name:           "space email",
			input:          `" "@example.org`,
			expectedResult: true,
			err:            nil,
		},
		{
			name:           "email without @",
			input:          "Abc.example.com",
			expectedResult: false,
			err:            fmt.Errorf("Abc.example.com is not valid email address, the format of email addresses is local-part@domain"),
		},
		{
			name:           "multiple @",
			input:          "A@b@c@example.com",
			expectedResult: false,
			err:            fmt.Errorf("an email address can't have multiple '@' characters"),
		},
		{
			name:           "Quote at domain",
			input:          `test@ex"ample.com`,
			expectedResult: false,
			err:            fmt.Errorf(`ex"ample.com is not a valid domain`),
		},
		{
			name:           "email start with @",
			input:          `@example.com`,
			expectedResult: false,
			err:            fmt.Errorf("email address can't start with '@'"),
		},
		{
			name:           "email end with @",
			input:          `test@`,
			expectedResult: false,
			err:            fmt.Errorf("domain part can't be empty"),
		},
		{
			name:           "exceed maximum domain length",
			input:          `test@abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz`,
			expectedResult: false,
			err:            fmt.Errorf("%s is longer than %d", "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz", MaxDomainLength),
		},
		{
			name:           "exceed maximum local part",
			input:          `abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz@example.com`,
			expectedResult: false,
			err:            fmt.Errorf("the length of local part should be less than %d", MaxLocalPart),
		},
		{
			name:           "email with escape",
			input:          `te\st@test.net`,
			expectedResult: false,
			err:            fmt.Errorf("fail to parse localPart of the email address"),
		},
		{
			name:           "customer/department=shipping@example.com",
			input:          `customer/department=shipping@example.com`,
			expectedResult: true,
			err:            nil,
		},
		{
			name:           "$A12345@example.com",
			input:          "$A12345@example.com",
			expectedResult: true,
			err:            nil,
		},
		{
			name:           "!def!xyz%abc@example.com",
			input:          "!def!xyz%abc@example.com",
			expectedResult: true,
			err:            nil,
		},
		{
			name:           "_somename@example.com",
			input:          "_somename@example.com",
			expectedResult: true,
			err:            nil,
		},
		{
			name:           "multiple dot and tag",
			input:          "disposable.style.email.with+symbol@example.com",
			expectedResult: true,
			err:            nil,
		},
		{
			name:           "simple and legit",
			input:          "simple@example.com",
			expectedResult: true,
			err:            nil,
		},
		{
			name:           "multiple hyphen",
			input:          "other.email-with-hyphen@example.com",
			expectedResult: true,
			err:            nil,
		},
		{
			name:           "hyphens",
			input:          "fully-qualified-domain@example.com",
			expectedResult: true,
			err:            nil,
		},
		{
			name:           "two tags",
			input:          "user.name+tag+sorting@example.com",
			expectedResult: true,
			err:            nil,
		},
		{
			name:           "single letter on local part",
			input:          "x@example.com",
			expectedResult: true,
			err:            nil,
		},
		{
			name:           "hyphen in domain",
			input:          "example-indeed@strange-example.com",
			expectedResult: true,
			err:            nil,
		},
	}
	for _, item := range cases {
		t.Run(item.name, func(st *testing.T) {
			r, err := Validate(item.input)
			if nil != err && nil == item.err {
				st.Errorf("we are not expecting error , however we got :%s", err.Error())
				st.FailNow()
			}
			if nil == err && nil != item.err {
				st.Errorf("we expecting err:%s,however we got nil", item.err.Error())
				st.FailNow()
			}
			if nil != err && nil != item.err {
				if err.Error() != item.err.Error() {
					st.Errorf("we are expecting err:%s, however we got :%s", item.err, err)
					st.FailNow()
				}
			}
			if r != item.expectedResult {
				st.Errorf("expected result is %t however we got %t", item.expectedResult, r)
				st.FailNow()
			}
		})
	}
}

func TestParseLocalPart(t *testing.T) {
	cases := []struct {
		name           string
		input          string
		expectedResult *localPart
		err            error
	}{
		{
			name:           "empty local part",
			input:          ``,
			expectedResult: nil,
			err:            fmt.Errorf("empty local part"),
		},
		{
			name:  "escaped at sign",
			input: `Abc\@def`,
			expectedResult: &localPart{
				localPartEmail: `Abc\@def`,
			},
			err: nil,
		},
		{
			name:  "multiple tags",
			input: `johnny+asdf1+asdf2`,
			expectedResult: &localPart{
				localPartEmail: "johnny",
				tags: []tag{
					{
						emailTags: "+asdf1+asdf2",
						start:     1,
						end:       6,
					},
					{
						emailTags: "+asdf1+asdf2",
						start:     7,
						end:       12,
					},
				},
			},
			err: nil,
		},
		{
			name:  "with a dot",
			input: "very.common",
			expectedResult: &localPart{
				localPartEmail: "very.common",
			},
			err: nil,
		},
		{
			name:           "start with dot",
			input:          ".asdf",
			expectedResult: nil,
			err:            fmt.Errorf(". can't be the start or end of local part"),
		},
		{
			name:           "end with dot",
			input:          "asdf.",
			expectedResult: nil,
			err:            fmt.Errorf(". can't be the start or end of local part"),
		},
		{
			name:           "with a double quote",
			input:          `asdf"d`,
			expectedResult: nil,
			err:            fmt.Errorf("\" is only valid escaped with baskslash"),
		},
		{
			name:  "with a escaped double quote",
			input: `asdf\"d`,
			expectedResult: &localPart{
				localPartEmail: `asdf\"d`,
			},
			err: nil,
		},
		{
			name:  "space in double quotetation",
			input: `" "`,
			expectedResult: &localPart{
				localPartEmail: `" "`,
			},
			err: nil,
		},
		{
			name:  "escaped double slash",
			input: `\\Blow`,
			expectedResult: &localPart{
				localPartEmail: `\\Blow`,
			},
			err: nil,
		},
		{
			name:  "at sign in quotation",
			input: `"abc@def"`,
			expectedResult: &localPart{
				localPartEmail: `"abc@def"`,
			},
			err: nil,
		},
		{
			name:  "space in quotation",
			input: `"Fred Bloggs"`,
			expectedResult: &localPart{
				localPartEmail: `"Fred Bloggs"`,
			},
			err: nil,
		},
		{
			name:           "space without quotation",
			input:          `Fred Bloggs`,
			expectedResult: nil,
			err:            fmt.Errorf("%c is only valid in quoted string or escaped", ' '),
		},
		{
			name:           "start with bracket but no end",
			input:          "(test",
			expectedResult: nil,
			err:            fmt.Errorf("( is only valid within quoted string or escaped"),
		},
		{
			name:  "comment at the begining",
			input: `(test)jo`,
			expectedResult: &localPart{
				localPartEmail:    "jo",
				comment:           "test",
				commentAtBegining: true,
			},
			err: nil,
		},
		{
			name:  "a",
			input: `a`,
			expectedResult: &localPart{
				localPartEmail: "a",
			},
			err: nil,
		},
		{
			name:           ")",
			input:          `)`,
			expectedResult: nil,
			err:            fmt.Errorf("%s is invalid in the local part of an email address", ")"),
		},
		{
			name:           "abc)d",
			input:          `abc)d`,
			expectedResult: nil,
			err:            fmt.Errorf(") is only valid within quoted string or escaped"),
		},
		{
			name:           "localpart with )(",
			input:          "test)wel(come",
			expectedResult: nil,
			err:            fmt.Errorf("invalid email address"),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(st *testing.T) {
			lp, err := parseLocalPart(c.input)
			if nil != err && c.err == nil {
				st.Errorf("we are not expecting error , however we got:%s", err)
				st.FailNow()
			}
			if nil == err && c.err != nil {
				st.Errorf("we are expecting err:%s, however we got nil", c.err)
				st.FailNow()
			}
			if c.err != nil && err != nil {
				if c.err.Error() != err.Error() {
					st.Errorf("we are expecting err:%s,however we got :%s", c.err, err)
					st.FailNow()
				}
			}
			if c.expectedResult == nil && lp != nil {
				st.Errorf("we expect result to be nil , however we got : %s", lp)
				st.FailNow()
			}
			if c.expectedResult != nil && lp == nil {
				st.Errorf("we are expecting %s, however we got nil", c.expectedResult)
				st.FailNow()
			}
			if nil != lp && nil != c.expectedResult {
				if lp.String() != c.expectedResult.String() {
					st.Errorf("we are expecting %s, however we got :%s", c.expectedResult, lp)
					st.FailNow()
				}
			}
		})
	}
}

func TestEquals(t *testing.T) {
	cases := []struct {
		name        string
		inputFirst  string
		inputSecond string
		expected    bool
	}{
		{
			name:        "empty vs empty",
			inputFirst:  ``,
			inputSecond: ``,
			expected:    false,
		},
		{
			name:        "empty vs none empty",
			inputFirst:  ``,
			inputSecond: `test@test.net`,
			expected:    false,
		},
		{
			name:        "comment vs without comments",
			inputFirst:  `test(explaintest)@test.net`,
			inputSecond: `test@test.net`,
			expected:    true,
		},
		{
			name:        "comment vs comments",
			inputFirst:  `test(explaintest)@test.net`,
			inputSecond: `(blalala)test@test.net`,
			expected:    true,
		},
		{
			name:        "tag",
			inputFirst:  `test(explaintest)@test.net`,
			inputSecond: `test+hello@test.net`,
			expected:    true,
		},
		{
			name:        "tag vs no tag",
			inputFirst:  `test@test.net`,
			inputSecond: `test+hello@test.net`,
			expected:    true,
		},
		{
			name:        "mixcase",
			inputFirst:  `tesT@test.net`,
			inputSecond: `TesT+hello@test.net`,
			expected:    true,
		},
		{
			name:        "tag vs multiple tags",
			inputFirst:  `test+hello+hello1@test.net`,
			inputSecond: `test+hello@test.net`,
			expected:    true,
		},
		{
			name:        "first email is invalid",
			inputFirst:  "invalid*@test@test.net",
			inputSecond: "test@test.net",
			expected:    false,
		},
		{
			name:        "second email is invalid",
			inputFirst:  "test@test.net",
			inputSecond: "invalid*@test@test.net",
			expected:    false,
		},
		{
			name:        "different domain",
			inputFirst:  "test@test.net",
			inputSecond: "test@gmail.com",
			expected:    false,
		},
		{
			name:        "differnt local part",
			inputFirst:  "test@test.net",
			inputSecond: "test111@test.net",
			expected:    false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(st *testing.T) {
			result := Equals(c.inputFirst, c.inputSecond)
			if c.expected != result {
				st.Errorf("we expect %t however we got %t", c.expected, result)
			}
		})
	}
}

func benchEmailAddress(input string, expectedResult bool, b *testing.B) {
	b.ReportAllocs()
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	for i := 0; i < b.N; i++ {
		if expectedResult != re.MatchString(input) {
			b.Fail()
		}
	}
}

func benchEmailAddressPkg(input string, expectedResult bool, b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		result, err := Validate(input)
		if nil != err {
			b.Fail()
		}
		if result != expectedResult {
			b.Fail()
		}
	}
}
func BenchmarkEmailAddresspkg_testattestdotnet(b *testing.B) {
	benchEmailAddressPkg("test@test.net", true, b)
}
func BenchmarkEmailAddress_testattestdotnet(b *testing.B) { benchEmailAddress("test@test.net", true, b) }
func BenchmarkEmailAddress_testattestdotnetWithSpecialChar(b *testing.B) {
	benchEmailAddress("te#!sdt@test.net", true, b)
}

func BenchmarkValidateEmailAddress(b *testing.B) {
	inputs := []struct {
		name           string
		input          string
		expectedResult bool
	}{
		{
			name:           "With Comment At the begining",
			input:          "(hello)test@test.net",
			expectedResult: true,
		},
		{
			name:           "With Comment At the end",
			input:          "test(hello)@test.net",
			expectedResult: true,
		},
		{
			name:           "With One Tag",
			input:          "test+test1@test.net",
			expectedResult: true,
		},
		{
			name:           "Email with special char",
			input:          "te#!sdt@test.net",
			expectedResult: true,
		},
		{
			name:           "test@test.net",
			input:          "test@test.net",
			expectedResult: true,
		},
		{
			name:           "email with quote",
			input:          "\"display\"@helloworld.net",
			expectedResult: true,
		},
	}
	for _, item := range inputs {
		b.Run(item.name, func(sb *testing.B) {
			benchEmailAddressPkg(item.input, item.expectedResult, sb)
		})
	}
}

func ExampleValidate() {
	e := "test@test.net"
	b, err := Validate(e)
	if nil != err {
		panic(err)
	}
	if b {
		fmt.Printf("%s is a legitimate email address", e)
	} else {
		fmt.Printf("%s is not a legitimate email", e)
	}
}

// TestGetTags unit tests
func TestGetTags(t *testing.T) {
	testCases := []struct {
		name           string
		input          string
		expectedResult []tag
	}{
		{
			name:           "empty tags",
			input:          "",
			expectedResult: nil,
		},
		{
			name:           "only plus",
			input:          "+",
			expectedResult: nil,
		},
		{
			name:  "start with plus",
			input: "+test+test1",
			expectedResult: []tag{
				{
					emailTags: "+test+test1",
					start:     1,
					end:       5,
				},
				{
					emailTags: "+test+test1",
					start:     6,
					end:       11,
				},
			},
		},
		{
			name:  "end with plus",
			input: "test+test1+",
			expectedResult: []tag{
				{
					emailTags: "test+test1+",
					start:     0,
					end:       4,
				},
				{
					emailTags: "test+test1+",
					start:     5,
					end:       10,
				},
			},
		},
		{
			name:  "only one tag",
			input: "+test",
			expectedResult: []tag{
				{
					emailTags: "+test",
					start:     1,
					end:       5,
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(st *testing.T) {
			tags := getTags(tc.input)
			if tc.expectedResult == nil && tags != nil {
				st.Errorf("we are not expecting anything but we got : %v", tags)
				return
			}
			if !reflect.DeepEqual(tags, tc.expectedResult) {
				st.Errorf("we are expecting %v , however we got : %v", tc.expectedResult, tags)
				return
			}
		})
	}
}
