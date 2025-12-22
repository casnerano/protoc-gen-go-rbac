package rbac

import (
	"path"
	"strings"
)

func CheckAccess(service *Service, fullMethod string, roles []string) bool {
	if method, exists := service.Methods[path.Base(fullMethod)]; exists {
		return hasRolesAccess(method.Rules, roles)
	}
	return hasRolesAccess(service.Rules, roles)
}

func hasRolesAccess(rules *Rules, roles []string) bool {
	switch rules.AccessLevel {
	case AccessLevelPublic:
		return true
	case AccessLevelPrivate:
		for _, role := range roles {
			role = strings.ToLower(role)
			for _, allowed := range rules.AllowedRoles {
				if role == allowed {
					return true
				}
			}
		}
		return false
	default:
		return false
	}
}
