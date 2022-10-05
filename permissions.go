package django

import (
	"net/http"

	"golang.org/x/exp/slices"
)

var (
	SAFE_METHODS = []string{
		http.MethodGet,
		http.MethodHead,
		http.MethodOptions,
	}
)

type PermissionClass interface {
	// HasPermission return true if permission is granted, false otherwise.
	HasPermission(request *Request) bool
}

type and struct {
	permissions []PermissionClass
}

func (p and) HasPermission(request *Request) bool {
	for _, permission := range p.permissions {
		if ok := permission.HasPermission(request); !ok {
			return false
		}
	}
	return true
}

func And(permissions ...PermissionClass) PermissionClass {
	return and{
		permissions: permissions,
	}
}

type or struct {
	permissions []PermissionClass
}

func (p or) HasPermission(request *Request) bool {
	for _, permission := range p.permissions {
		if ok := permission.HasPermission(request); ok {
			return true
		}
	}
	return false
}

func Or(permissions ...PermissionClass) PermissionClass {
	return or{
		permissions: permissions,
	}
}

// BasePermission a base struct from which all permission classes should be composed.
type BasePermission struct{}

func (p BasePermission) HasPermission(request *Request) bool {
	return true
}

type AllowAny struct {
	BasePermission
}

func (p AllowAny) HasPermission(request *Request) bool {
	return true
}

type IsAuthenticated struct {
	BasePermission
}

func (p IsAuthenticated) HasPermission(request *Request) bool {
	return request.User != nil && request.User.IsAuthenticated()
}

type IsAdminUser struct {
	BasePermission
}

func (p IsAdminUser) HasPermission(request *Request) bool {
	return request.User != nil && request.User.IsAdmin()
}

type IsAuthenticatedOrReadOnly struct {
	BasePermission
}

func (p IsAuthenticatedOrReadOnly) HasPermission(request *Request) bool {
	return slices.Contains(SAFE_METHODS, request.Method) || (request.User != nil && request.User.IsAuthenticated())
}
