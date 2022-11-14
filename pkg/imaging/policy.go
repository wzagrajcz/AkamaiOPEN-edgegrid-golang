package imaging

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/wzagrajcz/AkamaiOPEN-edgegrid-golang/v3/pkg/edgegriderr"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Code in this package is using autogenerated code located in the policy.gen.go file.
// Generator code is located in `dxe-tools` repo, and it is using OpenApi schema from
// https://git.source.akamai.com/users/eleclair/repos/terraform/browse/docs/schemas

type (
	// Policies is an Image and Video Manager API interface for Policy
	//
	// See: https://techdocs.akamai.com/ivm/reference/api
	Policies interface {
		// ListPolicies lists all Policies for the given network and an account
		//
		// See: https://techdocs.akamai.com/ivm/reference/get-policies
		ListPolicies(context.Context, ListPoliciesRequest) (*ListPoliciesResponse, error)

		// GetPolicy gets specific policy by PolicyID
		GetPolicy(context.Context, GetPolicyRequest) (PolicyOutput, error)

		// UpsertPolicy creates or updates the configuration for a policy
		UpsertPolicy(context.Context, UpsertPolicyRequest) (*PolicyResponse, error)

		// DeletePolicy deletes a policy
		DeletePolicy(context.Context, DeletePolicyRequest) (*PolicyResponse, error)

		// GetPolicyHistory retrieves history of changes for a policy
		GetPolicyHistory(context.Context, GetPolicyHistoryRequest) (*GetPolicyHistoryResponse, error)

		// RollbackPolicy reverts a policy to its previous version and deploys it to the network
		RollbackPolicy(ctx context.Context, request RollbackPolicyRequest) (*PolicyResponse, error)
	}

	// ListPoliciesRequest describes the parameters of the ListPolicies request
	ListPoliciesRequest struct {
		Network     PolicyNetwork
		ContractID  string
		PolicySetID string
	}

	// ListPoliciesResponse is a response returned by ListPolicies operations
	ListPoliciesResponse struct {
		ItemKind   string        `json:"itemKind"`
		Items      PolicyOutputs `json:"items"`
		TotalItems int           `json:"totalItems"`
	}

	// GetPolicyRequest describes the parameters of the GetPolicy request
	GetPolicyRequest policyRequest
	// DeletePolicyRequest describes the parameters of the DeletePolicy request
	DeletePolicyRequest policyRequest
	// GetPolicyHistoryRequest describes the parameters of the GetHistoryPolicy request
	GetPolicyHistoryRequest policyRequest
	// RollbackPolicyRequest describes the parameters of the RollbackPolicy request
	RollbackPolicyRequest policyRequest

	// policyRequest describes the parameters of the various policy requests
	policyRequest struct {
		PolicyID    string
		Network     PolicyNetwork
		ContractID  string
		PolicySetID string
	}

	// UpsertPolicyRequest describes the parameters of the UpsertPolicy request
	UpsertPolicyRequest struct {
		PolicyID    string
		Network     PolicyNetwork
		ContractID  string
		PolicySetID string
		PolicyInput
	}

	// PolicyResponse describes response of the UpsertPolicy, DeletePolicy and RollbackPolicy responses
	PolicyResponse struct {
		Description        string `json:"description"`
		ID                 string `json:"id"`
		OperationPerformed string `json:"operationPerformed"`
	}

	// GetPolicyHistoryResponse describes the parameters of the GetPolicyHistory response
	GetPolicyHistoryResponse struct {
		ItemKind   string              `json:"itemKind"`
		TotalItems int                 `json:"totalItems"`
		Items      []PolicyHistoryItem `json:"items"`
	}

	// PolicyHistoryItem describes items of the history for policy
	PolicyHistoryItem struct {
		ID          string `json:"id"`
		DateCreated string `json:"dateCreated"`
		Policy      string `json:"policy"`
		Action      string `json:"action"`
		User        string `json:"user"`
		Version     int    `json:"version"`
	}

	// PolicyOutput is implemented by PolicyOutput types (image and video)
	PolicyOutput interface {
		policyOutputType() string
	}

	// PolicyInput is implemented by PolicyInput types (image and video)
	PolicyInput interface {
		policyInputType() string
	}

	// PolicyOutputs is an array of PolicyOutput types (image and video)
	PolicyOutputs []PolicyOutput

	// PolicyNetwork represents the network where policy set is stored
	PolicyNetwork string

	// PolicyInputImage Specifies details for each policy, such as transformations to apply and variations in image size and formats
	PolicyInputImage struct {
		// Breakpoints The breakpoint widths (in pixels) to use to create derivative images/videos
		Breakpoints *Breakpoints `json:"breakpoints,omitempty"`
		// Hosts Hosts that are allowed for image/video URLs within transformations or variables
		Hosts []string `json:"hosts,omitempty"`
		// Output Dictates the output quality (either `quality` or `perceptualQuality`) and formats that are created for each resized image If unspecified, image formats are created to support all browsers at the default quality level (`85`), which includes formats such as WEBP, JPEG2000 and JPEG-XR for specific browsers
		Output *OutputImage `json:"output,omitempty"`
		// PostBreakpointTransformations Post-processing Transformations are applied to the image after image and quality settings have been applied
		PostBreakpointTransformations PostBreakpointTransformations `json:"postBreakpointTransformations,omitempty"`
		// RolloutDuration The amount of time in seconds that the policy takes to rollout. During the rollout an increasing proportion of images/videos will begin to use the new policy instead of the cached images/videos from the previous version
		RolloutDuration *int `json:"rolloutDuration,omitempty"`
		// Transformations Set of image transformations to apply to the source image. If unspecified, no operations are performed
		Transformations Transformations `json:"transformations,omitempty"`
		// Variables Declares variables for use within the policy. Any variable declared here can be invoked throughout transformations as a [Variable](#variable) object, so that you don't have to specify values separately You can also pass in these variable names and values dynamically as query parameters in the image's request URL
		Variables []Variable `json:"variables,omitempty"`
	}

	// PolicyInputVideo Specifies details for each policy such as video size
	PolicyInputVideo struct {
		// Breakpoints The breakpoint widths (in pixels) to use to create derivative images/videos
		Breakpoints *Breakpoints `json:"breakpoints,omitempty"`
		// Hosts Hosts that are allowed for image/video URLs within transformations or variables
		Hosts []string `json:"hosts,omitempty"`
		// Output Dictates the output quality that are created for each resized video
		Output *OutputVideo `json:"output,omitempty"`
		// RolloutDuration The amount of time in seconds that the policy takes to rollout. During the rollout an increasing proportion of images/videos will begin to use the new policy instead of the cached images/videos from the previous version
		RolloutDuration *int `json:"rolloutDuration,omitempty"`
		// Variables Declares variables for use within the policy. Any variable declared here can be invoked throughout transformations as a [Variable](#variable) object, so that you don't have to specify values separately You can also pass in these variable names and values dynamically as query parameters in the image's request URL
		Variables []Variable `json:"variables,omitempty"`
	}
)

const (
	// PolicyNetworkStaging represents staging network
	PolicyNetworkStaging PolicyNetwork = "staging"
	// PolicyNetworkProduction represents production network
	PolicyNetworkProduction PolicyNetwork = "production"
)

var (
	// ErrUnmarshalPolicyOutputList represents an error while unmarshalling transformation list
	ErrUnmarshalPolicyOutputList = errors.New("unmarshalling policy output list")

	// ErrListPolicies is returned when ListPolicies fails
	ErrListPolicies = errors.New("list policies")

	// ErrGetPolicy is returned when GetPolicy fails
	ErrGetPolicy = errors.New("get policy")

	// ErrUpsertPolicy is returned when UpsertPolicy fails
	ErrUpsertPolicy = errors.New("upsert policy")

	// ErrDeletePolicy is returned when DeletePolicy fails
	ErrDeletePolicy = errors.New("delete policy")

	// ErrGetPolicyHistory is returned when GetPolicyHistory fails
	ErrGetPolicyHistory = errors.New("get policy history")

	// ErrRollbackPolicy is returned when RollbackPolicy fails
	ErrRollbackPolicy = errors.New("rollback policy")
)

func (*PolicyOutputImage) policyOutputType() string {
	return "Image"
}

func (*PolicyOutputVideo) policyOutputType() string {
	return "Video"
}

func (*PolicyInputImage) policyInputType() string {
	return "Image"
}

func (*PolicyInputVideo) policyInputType() string {
	return "Video"
}

// Validate validates PolicyInputImage
func (p *PolicyInputImage) Validate() error {
	return validation.Errors{
		"Breakpoints":                   validation.Validate(p.Breakpoints),
		"Hosts":                         validation.Validate(p.Hosts, validation.Each()),
		"Output":                        validation.Validate(p.Output),
		"PostBreakpointTransformations": validation.Validate(p.PostBreakpointTransformations),
		"RolloutDuration": validation.Validate(p.RolloutDuration,
			validation.Min(3600),
			validation.Max(604800),
		),
		"Transformations": validation.Validate(p.Transformations),
		"Variables":       validation.Validate(p.Variables, validation.Each()),
	}.Filter()
}

// Validate validates PolicyInputVideo
func (p *PolicyInputVideo) Validate() error {
	return validation.Errors{
		"Breakpoints": validation.Validate(p.Breakpoints),
		"Hosts":       validation.Validate(p.Hosts, validation.Each()),
		"Output":      validation.Validate(p.Output),
		"RolloutDuration": validation.Validate(p.RolloutDuration,
			validation.Min(3600),
			validation.Max(604800),
		),
		"Variables": validation.Validate(p.Variables, validation.Each()),
	}.Filter()
}

var policyOutputHandlers = map[bool]func() PolicyOutput{
	false: func() PolicyOutput { return &PolicyOutputImage{} },
	true:  func() PolicyOutput { return &PolicyOutputVideo{} },
}

// UnmarshalJSON is a custom unmarshaler used to decode a slice of PolicyOutput interfaces
func (po *PolicyOutputs) UnmarshalJSON(in []byte) error {
	data := make([]map[string]interface{}, 0)
	if err := json.Unmarshal(in, &data); err != nil {
		return fmt.Errorf("%w: %s", ErrUnmarshalPolicyOutputList, err)
	}
	for _, policyOutput := range data {
		p, err := unmarshallPolicyOutput(policyOutput)
		if err != nil {
			return err
		}
		*po = append(*po, p)
	}
	return nil
}

func unmarshallPolicyOutput(policyOutput map[string]interface{}) (PolicyOutput, error) {
	video, ok := policyOutput["video"]
	if !ok {
		return nil, fmt.Errorf("%w: policyOutput should contain 'video' field", ErrUnmarshalPolicyOutputList)
	}
	isVideo, ok := video.(bool)
	if !ok {
		return nil, fmt.Errorf("%w: 'video' field on policyOutput entry should be a boolean", ErrUnmarshalPolicyOutputList)
	}

	bytes, err := json.Marshal(policyOutput)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrUnmarshalPolicyOutputList, err)
	}

	indicatedPolicyOutputType, ok := policyOutputHandlers[isVideo]
	if !ok {
		return nil, fmt.Errorf("%w: unsupported policyOutput type: %v", ErrUnmarshalPolicyOutputList, isVideo)
	}
	ipt := indicatedPolicyOutputType()
	err = json.Unmarshal(bytes, ipt)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrUnmarshalPolicyOutputList, err)
	}
	return ipt, nil
}

// Validate validates ListPoliciesRequest
func (v ListPoliciesRequest) Validate() error {
	errs := validation.Errors{
		"ContractID":  validation.Validate(v.ContractID, validation.Required),
		"PolicySetID": validation.Validate(v.PolicySetID, validation.Required),
		"Network": validation.Validate(v.Network, validation.Required, validation.In(PolicyNetworkStaging, PolicyNetworkProduction).
			Error(fmt.Sprintf("network has to be '%s', '%s'", PolicyNetworkStaging, PolicyNetworkProduction))),
	}
	return edgegriderr.ParseValidationErrors(errs)
}

// Validate validates GetPolicyRequest
func (v GetPolicyRequest) Validate() error {
	errs := validation.Errors{
		"PolicyID":    validation.Validate(v.PolicyID, validation.Required),
		"ContractID":  validation.Validate(v.ContractID, validation.Required),
		"PolicySetID": validation.Validate(v.PolicySetID, validation.Required),
		"Network": validation.Validate(v.Network, validation.Required, validation.In(PolicyNetworkStaging, PolicyNetworkProduction).
			Error(fmt.Sprintf("network has to be '%s', '%s'", PolicyNetworkStaging, PolicyNetworkProduction))),
	}
	return edgegriderr.ParseValidationErrors(errs)
}

// Validate validates UpsertPolicyRequest
func (v UpsertPolicyRequest) Validate() error {
	errs := validation.Errors{
		"PolicyID":    validation.Validate(v.PolicyID, validation.Required),
		"ContractID":  validation.Validate(v.ContractID, validation.Required),
		"PolicySetID": validation.Validate(v.PolicySetID, validation.Required),
		"Network": validation.Validate(v.Network, validation.Required, validation.In(PolicyNetworkStaging, PolicyNetworkProduction).
			Error(fmt.Sprintf("network has to be '%s', '%s'", PolicyNetworkStaging, PolicyNetworkProduction))),
		"Policy": validation.Validate(v.PolicyInput, validation.Required),
		//Validate, Policy Input
	}
	return edgegriderr.ParseValidationErrors(errs)
}

// Validate validates DeletePolicyRequest
func (v DeletePolicyRequest) Validate() error {
	errs := validation.Errors{
		"PolicyID":    validation.Validate(v.PolicyID, validation.Required),
		"ContractID":  validation.Validate(v.ContractID, validation.Required),
		"PolicySetID": validation.Validate(v.PolicySetID, validation.Required),
		"Network": validation.Validate(v.Network, validation.Required, validation.In(PolicyNetworkStaging, PolicyNetworkProduction).
			Error(fmt.Sprintf("network has to be '%s', '%s'", PolicyNetworkStaging, PolicyNetworkProduction))),
	}
	return edgegriderr.ParseValidationErrors(errs)
}

// Validate validates GetPolicyHistoryRequest
func (v GetPolicyHistoryRequest) Validate() error {
	errs := validation.Errors{
		"PolicyID":    validation.Validate(v.PolicyID, validation.Required),
		"ContractID":  validation.Validate(v.ContractID, validation.Required),
		"PolicySetID": validation.Validate(v.PolicySetID, validation.Required),
		"Network": validation.Validate(v.Network, validation.Required, validation.In(PolicyNetworkStaging, PolicyNetworkProduction).
			Error(fmt.Sprintf("network has to be '%s', '%s'", PolicyNetworkStaging, PolicyNetworkProduction))),
	}
	return edgegriderr.ParseValidationErrors(errs)
}

// Validate validates RollbackPolicyRequest
func (v RollbackPolicyRequest) Validate() error {
	errs := validation.Errors{
		"PolicyID":    validation.Validate(v.PolicyID, validation.Required),
		"ContractID":  validation.Validate(v.ContractID, validation.Required),
		"PolicySetID": validation.Validate(v.PolicySetID, validation.Required),
		"Network": validation.Validate(v.Network, validation.Required, validation.In(PolicyNetworkStaging, PolicyNetworkProduction).
			Error(fmt.Sprintf("network has to be '%s', '%s'", PolicyNetworkStaging, PolicyNetworkProduction))),
	}
	return edgegriderr.ParseValidationErrors(errs)
}

func (i *imaging) ListPolicies(ctx context.Context, params ListPoliciesRequest) (*ListPoliciesResponse, error) {
	logger := i.Log(ctx)
	logger.Debug("ListPolicies")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrListPolicies, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/imaging/v2/network/%s/policies/", params.Network)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrListPolicies, err)
	}

	req.Header.Set("Contract", params.ContractID)
	req.Header.Set("Policy-Set", params.PolicySetID)

	var result ListPoliciesResponse
	resp, err := i.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrListPolicies, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrListPolicies, i.Error(resp))
	}

	return &result, nil
}

func (i *imaging) GetPolicy(ctx context.Context, params GetPolicyRequest) (PolicyOutput, error) {
	logger := i.Log(ctx)
	logger.Debug("GetPolicy")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrGetPolicy, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/imaging/v2/network/%s/policies/%s", params.Network, params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetPolicy, err)
	}

	req.Header.Set("Contract", params.ContractID)
	req.Header.Set("Policy-Set", params.PolicySetID)

	var result map[string]interface{}
	resp, err := i.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetPolicy, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetPolicy, i.Error(resp))
	}

	policyOutput, err := unmarshallPolicyOutput(result)
	if err != nil {
		return nil, err
	}

	return policyOutput, nil
}

func (i *imaging) UpsertPolicy(ctx context.Context, params UpsertPolicyRequest) (*PolicyResponse, error) {
	logger := i.Log(ctx)
	logger.Debug("UpsertPolicy")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrUpsertPolicy, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/imaging/v2/network/%s/policies/%s", params.Network, params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrUpsertPolicy, err)
	}

	req.Header.Set("Contract", params.ContractID)
	req.Header.Set("Policy-Set", params.PolicySetID)

	var result PolicyResponse
	resp, err := i.Exec(req, &result, params.PolicyInput)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrUpsertPolicy, err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%s: %w", ErrUpsertPolicy, i.Error(resp))
	}

	return &result, nil
}

func (i *imaging) DeletePolicy(ctx context.Context, params DeletePolicyRequest) (*PolicyResponse, error) {
	logger := i.Log(ctx)
	logger.Debug("DeletePolicy")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrDeletePolicy, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/imaging/v2/network/%s/policies/%s", params.Network, params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrDeletePolicy, err)
	}

	req.Header.Set("Contract", params.ContractID)
	req.Header.Set("Policy-Set", params.PolicySetID)

	var result PolicyResponse
	resp, err := i.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrDeletePolicy, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrDeletePolicy, i.Error(resp))
	}

	return &result, nil
}

func (i *imaging) GetPolicyHistory(ctx context.Context, params GetPolicyHistoryRequest) (*GetPolicyHistoryResponse, error) {
	logger := i.Log(ctx)
	logger.Debug("GetPolicyHistory")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrGetPolicyHistory, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/imaging/v2/network/%s/policies/history/%s", params.Network, params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrGetPolicyHistory, err)
	}

	req.Header.Set("Contract", params.ContractID)
	req.Header.Set("Policy-Set", params.PolicySetID)

	var result GetPolicyHistoryResponse
	resp, err := i.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrGetPolicyHistory, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrGetPolicyHistory, i.Error(resp))
	}

	return &result, nil
}

func (i *imaging) RollbackPolicy(ctx context.Context, params RollbackPolicyRequest) (*PolicyResponse, error) {
	logger := i.Log(ctx)
	logger.Debug("RollbackPolicy")

	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("%s: %w:\n%s", ErrRollbackPolicy, ErrStructValidation, err)
	}

	uri := fmt.Sprintf("/imaging/v2/network/%s/policies/rollback/%s", params.Network, params.PolicyID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create request: %s", ErrRollbackPolicy, err)
	}

	req.Header.Set("Contract", params.ContractID)
	req.Header.Set("Policy-Set", params.PolicySetID)

	var result PolicyResponse
	resp, err := i.Exec(req, &result)
	if err != nil {
		return nil, fmt.Errorf("%w: request failed: %s", ErrRollbackPolicy, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", ErrRollbackPolicy, i.Error(resp))
	}

	return &result, nil
}
