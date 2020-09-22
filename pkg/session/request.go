package session

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
)

var (
	ErrInvalidArgument = errors.New("invalid arguments provided")
	ErrMarshaling      = errors.New("marshaling input")
	ErrUnmarshaling    = errors.New("unmarshaling output")
)

// Exec will sign and execute the request using the client edgegrid.Config
func (s *session) Exec(r *http.Request, out interface{}, in ...interface{}) (*http.Response, error) {
	if len(in) > 1 {
		return nil, fmt.Errorf("%w: %s", ErrInvalidArgument, "'in' argument must have 0 or 1 value")
	}
	log := s.Log(r.Context())

	// Apply any context header overrides
	if o, ok := r.Context().Value(contextOptionKey).(*contextOptions); ok {
		for k, v := range o.header {
			r.Header[k] = v
		}
	}

	if r.UserAgent() == "" {
		r.Header.Set("User-Agent", s.userAgent)
	}

	if r.Header.Get("Content-Type") == "" {
		r.Header.Set("Content-Type", "application/json")
	}

	if r.URL.Scheme == "" {
		r.URL.Scheme = "https"
	}

	if r.URL.Host == "" {
		r.URL.Host = s.config.Host
	}

	if len(in) > 0 {
		data, err := json.Marshal(in[0])
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrMarshaling, err)
		}

		r.Body = ioutil.NopCloser(bytes.NewBuffer(data))
	}

	if err := s.Sign(r); err != nil {
		return nil, err
	}

	if s.trace {
		data, err := httputil.DumpRequestOut(r, true)
		if err != nil {
			log.WithError(err).Error("Failed to dump request")
		} else {
			log.Debug(string(data))
		}
	}

	resp, err := s.client.Do(r)
	if err != nil {
		return nil, err
	}

	if s.trace {
		data, err := httputil.DumpResponse(resp, true)
		if err != nil {
			log.WithError(err).Error("Failed to dump response")
		} else {
			log.Debug(string(data))
		}
	}

	if out != nil {
		data, err := ioutil.ReadAll(resp.Body)
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(data))
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(data, out); err != nil {
			return nil, fmt.Errorf("%w: %s", ErrUnmarshaling, err)
		}
	}

	return resp, nil
}

// Sign will only sign a request
func (s *session) Sign(r *http.Request) error {
	s.config.SignRequest(r)
	return nil
}
