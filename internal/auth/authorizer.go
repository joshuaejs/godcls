// authorizer is a simple access-list implementations

package auth

import (
	"fmt"

	"github.com/casbin/casbin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// New instantiates an enforcer
func New(model, policy string) *Authorizer {
	enforcer := casbin.NewEnforcer(model, policy)
	return &Authorizer{
		enforcer: enforcer,
	}
}

// Authorizer type wraps a single method, Authorize, around Casbin's Enforce
type Authorizer struct {
	enforcer *casbin.Enforcer
}

// Authorize defers to Casbin's Enforce function
func (a *Authorizer) Authorize(subject, object, action string) error {
	if !a.enforcer.Enforce(subject, object, action) {
		msg := fmt.Sprintf("%s not permitted to %s to %s", subject, object, action)
		st := status.New(codes.PermissionDenied, msg)
		return st.Err()
	}
	return nil
}
