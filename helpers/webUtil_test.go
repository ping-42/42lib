package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testCases = []struct {
	name        string
	inputUrl    string
	expectedUrl string
	expectedErr error
}{
	{
		name:        "should work with scheme",
		inputUrl:    "https://google.com",
		expectedUrl: "google.com",
		expectedErr: nil,
	},
	{
		name:        "should work with www",
		inputUrl:    "https://www.google.com",
		expectedUrl: "www.google.com",
		expectedErr: nil,
	},
	{
		name:        "should work with path",
		inputUrl:    "https://www.google.com/somePath?someQ=1",
		expectedUrl: "www.google.com",
		expectedErr: nil,
	},
}

func TestDomainExtraction(t *testing.T) {
	t.Parallel()

	for _, tc := range testCases {
		t.Run(tc.name, func(t2 *testing.T) {
			actual, actualErr := ExtractDomainFromUrl(tc.inputUrl)
			assert.Equal(t2, tc.expectedErr, actualErr)
			assert.Equal(t2, tc.expectedUrl, actual)
		})
	}
}
