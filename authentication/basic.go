package authentication

import (
	"fmt"

	django "github.com/malijoe/djanGo-unchained"
)

type BasicAuthentication struct {
	Realm    string
	provider UserProvider
}

func (a *BasicAuthentication) Authenticate(request *django.Request) error {
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
