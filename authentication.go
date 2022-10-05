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
	// Authenticate
	// authenticates the request
	Authenticate(*Request) error
	// AuthenticateHeader
	// returns a string to be used as the value of the `WWW-Authenticate`
	// header in a `
	AuthenticateHeader() string
}

type UserProvider interface {
	GetUser(id uint) (User, error)
	FindUser(username string) (User, error)
}

var DEFAULT_AUTHENTICATION_CLASSES = []AuthenticationClass{
	&BasicAuthentication{},
	&SessionAuthentication{},
}

// BaseAuthentication a minimal implementation of the AuthenticationClass interface
// that all AuthenticationClass implementations should embed
type BaseAuthentication struct{}

func (b BaseAuthentication) Authenticate(_ *Request) error {
	return nil
}

func (b BaseAuthentication) AuthenticateHeader() string {
	return ""
}

type BasicAuthentication struct {
	BaseAuthentication
	realm    string
	provider UserProvider
}

func NewBasicAuthentication(provider UserProvider, realm ...string) AuthenticationClass {
	auth := BasicAuthentication{
		provider: provider,
		realm:    "api",
	}
	if len(realm) > 0 {
		auth.realm = realm[0]
	}
	return &auth
}

func (a *BasicAuthentication) Authenticate(request *Request) error {
	username, password, ok := request.BasicAuth()
	if ok {
		u, err := a.provider.FindUser(username)
		if err != nil {
			return err
		}
		if err = u.ValidatePassword(password); err != nil {
			return err
		}
		request.User = u
		return nil
	}
	return fmt.Errorf("%w: missing or mal-formed authentication header", ErrorCouldNotAuthenticate)
}

func (a *BasicAuthentication) AuthenticateHeader() string {
	return fmt.Sprintf("Basic realm=%s", a.realm)
}

type SessionAuthentication struct {
	BaseAuthentication
	manager  *scs.SessionManager
	provider UserProvider
}

func NewSessionAuthentication(provider UserProvider, manager *scs.SessionManager) AuthenticationClass {
	auth := SessionAuthentication{
		provider: provider,
		manager:  manager,
	}
	return &auth
}

func (a *SessionAuthentication) Authenticate(request *Request) error {
	if ctxId := a.manager.Get(request.Context(), "userID"); ctxId != nil {
		id, ok := ctxId.(uint)
		if !ok {
			panic("user id value in context is not of type uint")
		}
		u, err := a.provider.GetUser(id)
		if err != nil {
			return err
		}
		request.User = u
		return nil
	}
	return fmt.Errorf("%w: no session found", ErrorCouldNotAuthenticate)
}

type TokenAuthentication struct {
	BaseAuthentication
	keyword  string
	provider UserProvider
}

func NewTokenAuthentication(provider UserProvider, keyword ...string) AuthenticationClass {
	auth := TokenAuthentication{
		provider: provider,
		keyword:  "Token",
	}
	if len(keyword) > 0 {
		auth.keyword = keyword[0]
	}
	return &auth
}

func (a *TokenAuthentication) Authenticate(request *Request) error {
	auth := request.Request.Header.Get("Authorization")
	if auth != "" {

	}
	return fmt.Errorf("%w: missing authorization header", ErrorCouldNotAuthenticate)
}

func (a *TokenAuthentication) AuthenticateHeader() string {
	return a.keyword
}
