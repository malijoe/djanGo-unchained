package permissions

type PermissionClass uint32

const (
	AllowAny PermissionClass = iota
	IsAuthenticated
	IsAdminUser
	IsAuthenticatedOrReadOnly
)
