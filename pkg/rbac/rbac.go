package rbac

import (
	"path"
	"strings"

	desc "github.com/casnerano/protoc-gen-go-rbac/proto"
)

type Service struct {
	Name    string
	Rules   *desc.Rules
	Methods map[string]*Method
}

type Method struct {
	Rules *desc.Rules
}

func CheckAccess(service *Service, fullMethod string, roles []string) bool {
	if method, exists := service.Methods[path.Base(fullMethod)]; exists {
		return hasRolesAccess(method.Rules, roles)
	}
	return hasRolesAccess(service.Rules, roles)
}

func hasRolesAccess(rules *desc.Rules, roles []string) bool {
	switch rules.AccessLevel {
	case desc.AccessLevel_ACCESS_LEVEL_PUBLIC:
		return true
	case desc.AccessLevel_ACCESS_LEVEL_PRIVATE:
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
