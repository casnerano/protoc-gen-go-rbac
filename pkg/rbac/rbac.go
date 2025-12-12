package rbac

import (
	"strings"

	rbac "github.com/casnerano/protoc-gen-go-rbac/proto"
)

type AccessRules struct {
	AccessLevel  rbac.AccessLevel
	AllowedRoles []string
}

type AccessService struct {
	Rules   *AccessRules
	Methods map[string]AccessMethod
}

type AccessMethod struct {
	Rules *AccessRules
}

func CheckRolesAccess(accessServices map[string]AccessService, fullMethod string, roles []string) bool {
	serviceName, methodName := extractServiceMethodNames(fullMethod)

	service, existServiceRules := accessServices[serviceName]
	if !existServiceRules {
		return false
	}

	var methodRules *AccessRules
	if method, existMethodRules := service.Methods[methodName]; existMethodRules {
		methodRules = method.Rules
	}

	mergedRules := mergeServiceMethodRules(service.Rules, methodRules)
	if mergedRules == nil {
		return false
	}

	return hasRolesAccess(mergedRules, roles)
}

func hasRolesAccess(rules *AccessRules, roles []string) bool {
	switch rules.AccessLevel {
	case rbac.AccessLevel_ACCESS_LEVEL_PUBLIC:
		return true
	case rbac.AccessLevel_ACCESS_LEVEL_PRIVATE:
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

func extractServiceMethodNames(fullMethod string) (serviceName string, methodName string) {
	parts := strings.Split(fullMethod, "/")
	if len(parts) != 3 {
		return
	}

	lastDotIndex := strings.LastIndex(parts[1], ".")
	if lastDotIndex < 0 {
		return
	}

	serviceName, methodName = parts[1][lastDotIndex+1:], parts[2]

	return
}

func mergeServiceMethodRules(serviceRules, methodRules *AccessRules) *AccessRules {
	if serviceRules == nil {
		return methodRules
	}

	if methodRules == nil {
		return serviceRules
	}

	merged := *methodRules
	if merged.AccessLevel == rbac.AccessLevel_ACCESS_LEVEL_UNKNOWN {
		merged.AccessLevel = serviceRules.AccessLevel
	}

	if len(merged.AllowedRoles) == 0 {
		merged.AllowedRoles = serviceRules.AllowedRoles
	}

	return &merged
}
