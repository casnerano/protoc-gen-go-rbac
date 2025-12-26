package interceptor

import (
	"context"

	"github.com/casnerano/protoc-gen-go-rbac/pkg/rbac"
)

type authenticatedEvaluator struct{}

func newAuthenticatedEvaluator() *authenticatedEvaluator {
	return &authenticatedEvaluator{}
}

func (e authenticatedEvaluator) Evaluate(_ context.Context, _ *rbac.Rules, authContext *AuthContext, _ string) (bool, error) {
	return authContext.Authenticated, nil
}
