package appsec

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppSec_ListReputationProfileAction(t *testing.T) {

	result := GetReputationProfileActionsResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestReputationProfileAction/ReputationProfileActions.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           GetReputationProfileActionsRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetReputationProfileActionsResponse
		withError        error
		headers          http.Header
	}{
		"200 OK": {
			params: GetReputationProfileActionsRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			responseStatus:   http.StatusOK,
			responseBody:     string(respData),
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/reputation-profiles",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetReputationProfileActionsRequest{
				ConfigID: 43253,
				Version:  15,
				PolicyID: "AAAA_81230",
			},
			headers:        http.Header{},
			responseStatus: http.StatusInternalServerError,
			responseBody: `
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching propertys",
    "status": 500
}`,
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/reputation-profiles",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching propertys",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.GetReputationProfileActions(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers),
				),
				test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test ReputationProfileAction
func TestAppSec_GetReputationProfileAction(t *testing.T) {

	result := GetReputationProfileActionResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestReputationProfileAction/ReputationProfileAction.json"))
	json.Unmarshal([]byte(respData), &result)

	tests := map[string]struct {
		params           GetReputationProfileActionRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *GetReputationProfileActionResponse
		withError        error
	}{
		"200 OK": {
			params: GetReputationProfileActionRequest{
				ConfigID:            43253,
				Version:             15,
				PolicyID:            "AAAA_81230",
				ReputationProfileID: 134644,
			},
			responseStatus:   http.StatusOK,
			responseBody:     respData,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/reputation-profiles/134644",
			expectedResponse: &result,
		},
		"500 internal server error": {
			params: GetReputationProfileActionRequest{
				ConfigID:            43253,
				Version:             15,
				PolicyID:            "AAAA_81230",
				ReputationProfileID: 134644,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error fetching match target"
}`),
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/reputation-profiles/134644",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error fetching match target",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.String())
				assert.Equal(t, http.MethodGet, r.Method)
				w.WriteHeader(test.responseStatus)
				_, err := w.Write([]byte(test.responseBody))
				assert.NoError(t, err)
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.GetReputationProfileAction(context.Background(), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}

// Test Update ReputationProfileAction.
func TestAppSec_UpdateReputationProfileAction(t *testing.T) {
	result := UpdateReputationProfileActionResponse{}

	respData := compactJSON(loadFixtureBytes("testdata/TestReputationProfileAction/ReputationProfileAction.json"))
	json.Unmarshal([]byte(respData), &result)

	req := UpdateReputationProfileActionRequest{}

	reqData := compactJSON(loadFixtureBytes("testdata/TestReputationProfileAction/ReputationProfileAction.json"))
	json.Unmarshal([]byte(reqData), &req)

	tests := map[string]struct {
		params           UpdateReputationProfileActionRequest
		responseStatus   int
		responseBody     string
		expectedPath     string
		expectedResponse *UpdateReputationProfileActionResponse
		withError        error
		headers          http.Header
	}{
		"200 Success": {
			params: UpdateReputationProfileActionRequest{
				ConfigID:            43253,
				Version:             15,
				PolicyID:            "AAAA_81230",
				ReputationProfileID: 134644,
			},
			headers: http.Header{
				"Content-Type": []string{"application/json;charset=UTF-8"},
			},
			responseStatus:   http.StatusCreated,
			responseBody:     respData,
			expectedResponse: &result,
			expectedPath:     "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/reputation-profiles/134644",
		},
		"500 internal server error": {
			params: UpdateReputationProfileActionRequest{
				ConfigID:            43253,
				Version:             15,
				PolicyID:            "AAAA_81230",
				ReputationProfileID: 134644,
			},
			responseStatus: http.StatusInternalServerError,
			responseBody: (`
{
    "type": "internal_error",
    "title": "Internal Server Error",
    "detail": "Error creating zone"
}`),
			expectedPath: "/appsec/v1/configs/43253/versions/15/security-policies/AAAA_81230/reputation-profiles/134644",
			withError: &Error{
				Type:       "internal_error",
				Title:      "Internal Server Error",
				Detail:     "Error creating zone",
				StatusCode: http.StatusInternalServerError,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPut, r.Method)
				w.WriteHeader(test.responseStatus)
				if len(test.responseBody) > 0 {
					_, err := w.Write([]byte(test.responseBody))
					assert.NoError(t, err)
				}
			}))
			client := mockAPIClient(t, mockServer)
			result, err := client.UpdateReputationProfileAction(
				session.ContextWithOptions(
					context.Background(),
					session.WithContextHeaders(test.headers)), test.params)
			if test.withError != nil {
				assert.True(t, errors.Is(err, test.withError), "want: %s; got: %s", test.withError, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectedResponse, result)
		})
	}
}