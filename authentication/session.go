package authentication

import (
	"fmt"

	"github.com/alexedwards/scs/v2"
	django "github.com/malijoe/djanGo-unchained"
)

type SessionAuthentication struct {
	manager  *scs.SessionManager
	provider UserProvider
}

func (a *SessionAuthentication) Authenticate(request *django.Request) error {
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
