package appsec

import (
	"context"
	"fmt"
	"net/http"
)

// ApiHostnameCoverage represents a collection of ApiHostnameCoverage
//
// See: ApiHostnameCoverage.GetApiHostnameCoverage()
// API Docs: // appsec v1
//
// https://developer.akamai.com/api/cloud_security/application_security/v1.html

type (
	// ApiHostnameCoverage  contains operations available on ApiHostnameCoverage  resource
	// See: // appsec v1
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getapihostnamecoverage
	ApiHostnameCoverage interface {
		GetApiHostnameCoverage(ctx context.Context, params GetApiHostnameCoverageRequest) (*GetApiHostnameCoverageResponse, error)
	}

	GetApiHostnameCoverageRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		Hostname string `json:"-"`
	}

	GetApiHostnameCoverageResponse struct {
		HostnameCoverage []struct {
			Configuration struct {
				ID      int    `json:"id"`
				Name    string `json:"name"`
				Version int    `json:"version"`
			} `json:"configuration"`
			Status         string   `json:"status"`
			HasMatchTarget bool     `json:"hasMatchTarget"`
			Hostname       string   `json:"hostname"`
			PolicyNames    []string `json:"policyNames"`
		} `json:"hostnameCoverage"`
	}
)

// Validate validates GetApiHostnameCoverageRequest
/*func (v GetApiHostnameCoverageRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}
*/
func (p *appsec) GetApiHostnameCoverage(ctx context.Context, params GetApiHostnameCoverageRequest) (*GetApiHostnameCoverageResponse, error) {
	/*	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
		}
	*/
	logger := p.Log(ctx)
	logger.Debug("GetApiHostnameCoverage")

	var rval GetApiHostnameCoverageResponse

	uri := 
		"/appsec/v1/hostname-coverage"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getapihostnamecoverage request: %w", err)
	}

	resp, err := p.Exec(req, &rval)
	if err != nil {
		return nil, fmt.Errorf("getapihostnamecoverage  request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &rval, nil

}