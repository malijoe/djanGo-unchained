package authentication

import django "github.com/malijoe/djanGo-unchained"

type Class interface {
	Authenticate(*django.Request) error
	Method() string
}

type UserProvider interface {
	GetUser(id uint) (*django.User, error)
	FindUser(username string) (*django.User, error)
}

var DEFAULT_AUTHENTICATION_CLASSES = []Class{
	&BasicAuthentication{},
	&SessionAuthentication{},
}

const (
//BasicAuthentication AuthenticationClass = iota
//TokenAuthentication
//SessionAuthentication
//RemoteUserAuthentication
)
