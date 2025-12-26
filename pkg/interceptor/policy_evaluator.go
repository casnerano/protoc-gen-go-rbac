package interceptor

import (
	"context"
	"fmt"

	"github.com/casnerano/protoc-gen-go-rbac/pkg/rbac"
)

type (
	Policy   func(ctx context.Context, authCtx *AuthContext, fullMethod string) (bool, error)
	Policies map[string]Policy
)

type policyEvaluator struct {
	policies Policies
}

func newPolicyEvaluator(policies Policies) *policyEvaluator {
	return &policyEvaluator{
		policies: policies,
	}
}

func (e policyEvaluator) Evaluate(ctx context.Context, rules *rbac.Rules, authContext *AuthContext, fullMethod string) (bool, error) {
	fn, ok := e.policies[*rules.PolicyName]
	if !ok {
		return false, fmt.Errorf("policy not found: %s", *rules.PolicyName)
	}
	return fn(ctx, authContext, fullMethod)
}
