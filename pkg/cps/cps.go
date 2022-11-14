package cps

import (
	"errors"

	"github.com/wzagrajcz/AkamaiOPEN-edgegrid-golang/v3/pkg/session"
)

var (
	// ErrStructValidation is returned returned when given struct validation failed
	ErrStructValidation = errors.New("struct validation")
)

type (
	// CPS is the cps api interface
	CPS interface {
		Enrollments
		ChangeOperations
		DVChallenges
		PreVerification
	}

	cps struct {
		session.Session
	}

	// Option defines a CPS option
	Option func(*cps)

	// ClientFunc is a cps client new method, this can used for mocking
	ClientFunc func(sess session.Session, opts ...Option) CPS
)

// Client returns a new cps Client instance with the specified controller
func Client(sess session.Session, opts ...Option) CPS {
	c := &cps{
		Session: sess,
	}

	for _, opt := range opts {
		opt(c)
	}
	return c
}
