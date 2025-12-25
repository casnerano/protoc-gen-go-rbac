package rbac

type AccessLevel int

const (
	AccessLevelUnknown = iota
	AccessLevelPublic
	AccessLevelAuthenticated
	AccessLevelPrivate
	AccessLevelPolicy
)

type Rules struct {
	AccessLevel  AccessLevel
	AllowedRoles []string
	PolicyName   *string
}

type Service struct {
	Name    string
	Rules   *Rules
	Methods map[string]*Method
}

type Method struct {
	Rules *Rules
}
