package emailaddress

import (
	"fmt"
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
			err:            fmt.Errorf("consective dot only valid inside quotation"),
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
			err:            fmt.Errorf(`" is invalid character in domain part`),
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

func TestParseEmail(t *testing.T) {
	input := []string{
		`Abc\@def@example.com`,
		`Fred\ Bloggs@example.com`,
		`\\Blow@example.com`,
		`"Abc@def"@example.com`,
		`"Fred Bloggs"@example.com`,
		`.@test.net`,
		`asdfasdf.@test.net`,
	}
	for _, item := range input {
		_, err := parseEmailAddress(item)
		if nil != err {
			panic(err)
		}
	}
}

func TestParseLocalPart(t *testing.T) {
	input := []string{
		`Abc\@def`,
		`johnny+asdf1+asdf2`,
		`very.common`,
		`.asdf`,
		`asdf.`,
		`asdf"d`,
		`" "`,
		`\\Blow`,
		`"abc@def"`,
		`"Fred Bloggs"`,
		`Fred Bloggs`,
		`(test`,
		`(test)jo`,
	}
	for _, item := range input {
		lp, err := parseLocalPart(item)
		if nil != err {
			t.Error(err)
		}
		fmt.Printf("%+v\n", lp)
	}
}
