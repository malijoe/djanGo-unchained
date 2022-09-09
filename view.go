package django

type View struct {
	AuthenticationClasses []AuthenticationClass
	PermissionClasses     []PermissionClass
	AllowedMethods        []string
}
