package interceptor

import (
	"context"

	"github.com/casnerano/protoc-gen-go-rbac/pkg/rbac"
)

type evaluatorOptions struct {
}

type evaluator interface {
	Evaluate(ctx context.Context, rules *rbac.Rules, authContext *AuthContext, fullMethod string) (bool, error)
}
