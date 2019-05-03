package emailaddress

import "testing"

func TestIsDomainName(t *testing.T) {
	cases := []struct {
		name           string
		input          string
		expectedResult bool
	}{
		{
			name:           "empty",
			input:          ``,
			expectedResult: false,
		},
		{
			name:           "test.net",
			input:          `test.net`,
			expectedResult: true,
		},
		{
			name:           "dot before dash",
			input:          `test.-net`,
			expectedResult: false,
		},
		{
			name:           "double quote",
			input:          `"test".net`,
			expectedResult: false,
		},
		{
			name:           "dotdot",
			input:          `test..net`,
			expectedResult: false,
		},
		{
			name:           "too long part",
			input:          `abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz.net`,
			expectedResult: false,
		},
		{
			name:           "too long part with dash",
			input:          `abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz-test.net`,
			expectedResult: false,
		},
		{
			name:           "start with dot",
			input:          ".test.net",
			expectedResult: false,
		},
		{
			name:           "with digit",
			input:          "t1est.net",
			expectedResult: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(st *testing.T) {
			result := IsDomainName(c.input)
			if result != c.expectedResult {
				st.Errorf("we expected : %t , however we got : %t", c.expectedResult, result)
			}
		})
	}
}
