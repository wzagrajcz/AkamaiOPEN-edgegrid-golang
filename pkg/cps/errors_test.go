package cps

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/wzagrajcz/AkamaiOPEN-edgegrid-golang/v3/pkg/session"
	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
)

func TestNewError(t *testing.T) {
	sess, err := session.New()
	require.NoError(t, err)

	req, err := http.NewRequestWithContext(
		context.TODO(),
		http.MethodHead,
		"/",
		nil)
	require.NoError(t, err)

	tests := map[string]struct {
		response *http.Response
		expected *Error
	}{
		"valid response, status code 500": {
			response: &http.Response{
				Status:     "Internal Server Error",
				StatusCode: http.StatusInternalServerError,
				Body: ioutil.NopCloser(strings.NewReader(
					`{"type":"a","title":"b","detail":"c"}`),
				),
				Request: req,
			},
			expected: &Error{
				Type:       "a",
				Title:      "b",
				Detail:     "c",
				StatusCode: http.StatusInternalServerError,
			},
		},
		"invalid response body, assign status code": {
			response: &http.Response{
				Status:     "Internal Server Error",
				StatusCode: http.StatusInternalServerError,
				Body: ioutil.NopCloser(strings.NewReader(
					`test`),
				),
				Request: req,
			},
			expected: &Error{
				Title:      "test",
				Detail:     "",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res := Client(sess).(*cps).Error(test.response)
			assert.Equal(t, test.expected, res)
		})
	}
}
