package interceptor

import (
	"context"

	"github.com/casnerano/protoc-gen-go-rbac/pkg/rbac"
)

type publicEvaluator struct{}

func newPublicEvaluator() *publicEvaluator {
	return &publicEvaluator{}
}

func (e publicEvaluator) Evaluate(_ context.Context, _ *rbac.Rules, _ *AuthContext, _ string) (bool, error) {
	return true, nil
}
