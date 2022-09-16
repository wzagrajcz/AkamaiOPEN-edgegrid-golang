package appsec

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// The EvalRule interface supports retrieving and modifying the rules available for
	// evaluation and their actions, or the action of a specific rule.
	//
	// https://developer.akamai.com/api/cloud_security/application_security/v1.html#evalrule
	EvalRule interface {
		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getevalrules
		GetEvalRules(ctx context.Context, params GetEvalRulesRequest) (*GetEvalRulesResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#getevalrule
		GetEvalRule(ctx context.Context, params GetEvalRuleRequest) (*GetEvalRuleResponse, error)

		// https://developer.akamai.com/api/cloud_security/application_security/v1.html#putevalrule
		UpdateEvalRule(ctx context.Context, params UpdateEvalRuleRequest) (*UpdateEvalRuleResponse, error)
	}

	// GetEvalRulesRequest is used to retrieve the rules available for evaluation and their actions.
	GetEvalRulesRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		RuleID   int    `json:"-"`
	}

	// GetEvalRulesResponse is returned from a call to GetEvalRules.
	GetEvalRulesResponse struct {
		Rules []struct {
			ID                 int                     `json:"id,omitempty"`
			Action             string                  `json:"action,omitempty"`
			ConditionException *RuleConditionException `json:"conditionException,omitempty"`
		} `json:"evalRuleActions,omitempty"`
	}

	// GetEvalRuleRequest is used to retrieve a rule available for evaluation and its action.
	GetEvalRuleRequest struct {
		ConfigID int    `json:"-"`
		Version  int    `json:"-"`
		PolicyID string `json:"-"`
		RuleID   int    `json:"ruleId"`
	}

	// GetEvalRuleResponse is returned from a call to GetEvalRule.
	GetEvalRuleResponse struct {
		Action             string                  `json:"action,omitempty"`
		ConditionException *RuleConditionException `json:"conditionException,omitempty"`
	}

	// UpdateEvalRuleRequest is used to modify a rule available for evaluation and its action.
	UpdateEvalRuleRequest struct {
		ConfigID       int             `json:"-"`
		Version        int             `json:"-"`
		PolicyID       string          `json:"-"`
		RuleID         int             `json:"-"`
		Action         string          `json:"action"`
		JsonPayloadRaw json.RawMessage `json:"conditionException,omitempty"`
	}

	// UpdateEvalRuleResponse is returned from a call to UpdateEvalRule.
	UpdateEvalRuleResponse struct {
		Action             string                  `json:"action,omitempty"`
		ConditionException *RuleConditionException `json:"conditionException,omitempty"`
	}
)

// IsEmptyConditionException checks whether the ConditionException field is empty.
func (r *GetEvalRuleResponse) IsEmptyConditionException() bool {
	return r.ConditionException == nil
}

// Validate validates a GetEvalRuleRequest.
func (v GetEvalRuleRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
		"RuleID":   validation.Validate(v.RuleID, validation.Required),
	}.Filter()
}

// Validate validates a GetEvalRulesRequest.
func (v GetEvalRulesRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
	}.Filter()
}

// Validate validates an UpdateEvalRuleRequest.
func (v UpdateEvalRuleRequest) Validate() error {
	return validation.Errors{
		"ConfigID": validation.Validate(v.ConfigID, validation.Required),
		"Version":  validation.Validate(v.Version, validation.Required),
		"PolicyID": validation.Validate(v.PolicyID, validation.Required),
		"RuleID":   validation.Validate(v.RuleID, validation.Required),
	}.Filter()
}

func (p *appsec) GetEvalRule(ctx context.Context, params GetEvalRuleRequest) (*GetEvalRuleResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetEvalRule")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	var result GetEvalRuleResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/eval-rules/%d?includeConditionException=true",
		params.ConfigID,
		params.Version,
		params.PolicyID,
		params.RuleID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetEvalRule request: %w", err)
	}

	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetEvalRule request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	return &result, nil

}

func (p *appsec) GetEvalRules(ctx context.Context, params GetEvalRulesRequest) (*GetEvalRulesResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("GetEvalRules")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	var result GetEvalRulesResponse
	var filteredResult GetEvalRulesResponse

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/eval-rules?includeConditionException=true",
		params.ConfigID,
		params.Version,
		params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GetEvalRules request: %w", err)
	}

	resp, err := p.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("GetEvalRules request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, p.Error(resp)
	}

	if params.RuleID != 0 {
		for _, val := range result.Rules {
			if val.ID == params.RuleID {
				filteredResult.Rules = append(filteredResult.Rules, val)
			}
		}
	} else {
		filteredResult = result
	}

	return &filteredResult, nil

}

func (p *appsec) UpdateEvalRule(ctx context.Context, params UpdateEvalRuleRequest) (*UpdateEvalRuleResponse, error) {
	logger := p.Log(ctx)
	logger.Debug("UpdateEvalRule")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrStructValidation, err.Error())
	}

	uri := fmt.Sprintf(
		"/appsec/v1/configs/%d/versions/%d/security-policies/%s/eval-rules/%d/action-condition-exception",
		params.ConfigID,
		params.Version,
		params.PolicyID,
		params.RuleID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create UpdateEvalRule request: %w", err)
	}

	var result UpdateEvalRuleResponse
	req.Header.Set("Content-Type", "application/json")
	resp, err := p.Exec(req, &result, params)
	if err != nil {
		return nil, fmt.Errorf("UpdateEvalRule request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, p.Error(resp)
	}

	return &result, nil
}