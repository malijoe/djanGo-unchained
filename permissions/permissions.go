package permissions

import (
	"net/http"

	django "github.com/malijoe/djanGo-unchained"
	"golang.org/x/exp/slices"
)

var (
	SAFE_METHODS = []string{
		http.MethodGet,
		http.MethodHead,
		http.MethodOptions,
	}
)

type Class interface {
	// HasPermission return true if permission is granted, false otherwise.
	HasPermission(request *django.Request) bool
}

// BasePermission a base struct from which all permission classes should be composed.
type BasePermission struct{}

func (p BasePermission) HasPermission(request *django.Request) bool {
	return true
}

type AllowAny struct {
	BasePermission
}

func (p AllowAny) HasPermission(request *django.Request) bool {
	return true
}

type IsAuthenticated struct {
	BasePermission
}

func (p IsAuthenticated) HasPermission(request *django.Request) bool {
	return request.User != nil && request.User.IsAuthenticated()
}

type IsAdminUser struct {
	BasePermission
}

func (p IsAdminUser) HasPermission(request *django.Request) bool {
	return request.User != nil && request.User.IsAdmin()
}

type IsAuthenticatedOrReadOnly struct {
	BasePermission
}

func (p IsAuthenticatedOrReadOnly) HasPermission(request *django.Request) bool {
	return slices.Contains(SAFE_METHODS, request.Method) || (request.User != nil && request.User.IsAuthenticated())
}
