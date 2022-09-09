package django

import (
	"errors"
	"fmt"

	"github.com/alexedwards/scs/v2"
)

var (
	ErrorCouldNotAuthenticate = errors.New("could not authenticate")
)

type AuthenticationClass interface {
	Authenticate(*Request) error
	Method() string
}

type UserProvider interface {
	GetUser(id uint) (*User, error)
	FindUser(username string) (*User, error)
}

var DEFAULT_AUTHENTICATION_CLASSES = []AuthenticationClass{
	&BasicAuthentication{},
	&SessionAuthentication{},
}

const (
//BasicAuthentication AuthenticationClass = iota
//TokenAuthentication
//SessionAuthentication
//RemoteUserAuthentication
)

type BasicAuthentication struct {
	Realm    string
	provider UserProvider
}

func (a *BasicAuthentication) Authenticate(request *Request) error {
	username, password, ok := request.BasicAuth()
	if ok {
		user, err := a.provider.FindUser(username)
		if err != nil {
			return err
		}
		if err = user.ValidatePassword(password); err != nil {
			return err
		}
		request.User = user
		return nil
	}
	return fmt.Errorf("%w: missing or mal-formed authentication header", ErrorCouldNotAuthenticate)
}

func (a *BasicAuthentication) Method() string {
	return fmt.Sprintf("Basic realm=%s", a.Realm)
}

type SessionAuthentication struct {
	manager  *scs.SessionManager
	provider UserProvider
}

func (a *SessionAuthentication) Authenticate(request *Request) error {
	if ctxId := a.manager.Get(request.Context(), "userID"); ctxId != nil {
		id, ok := ctxId.(uint)
		if !ok {
			panic("user id value in context is not of type uint")
		}
		user, err := a.provider.GetUser(id)
		if err != nil {
			return err
		}
		request.User = user
		return nil
	}
	return fmt.Errorf("%w: no session found", ErrorCouldNotAuthenticate)
}

func (a *SessionAuthentication) Method() string {
	return "Session"
}

type TokenAuthentication struct {
	provider UserProvider
}

func (a *TokenAuthentication) Authenticate(request *Request) error {
	auth := request.Request.Header.Get("Authorization")
	if auth != "" {

	}
	return fmt.Errorf("%w: missing authorization header", ErrorCouldNotAuthenticate)
}
