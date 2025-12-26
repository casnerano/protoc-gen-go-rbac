package interceptor

import (
	"context"

	"github.com/casnerano/protoc-gen-go-rbac/pkg/rbac"
)

type privateEvaluator struct{}

func newPrivateEvaluator() *privateEvaluator {
	return &privateEvaluator{}
}

func (e privateEvaluator) Evaluate(_ context.Context, rules *rbac.Rules, authContext *AuthContext, _ string) (bool, error) {
	if !authContext.Authenticated {
		return false, nil
	}

	for _, authRole := range authContext.Roles {
		for _, allowedRole := range rules.AllowedRoles {
			if authRole == allowedRole {
				return true, nil
			}
		}
	}

	return false, nil
}
