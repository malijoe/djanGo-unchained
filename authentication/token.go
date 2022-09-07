package authentication

import (
	"fmt"

	django "github.com/malijoe/djanGo-unchained"
)

type TokenAuthentication struct {
	provider UserProvider
}

func (a *TokenAuthentication) Authenticate(request *django.Request) error {
	auth := request.Request.Header.Get("Authorization")
	if auth != "" {

	}
	return fmt.Errorf("%w: missing authorization header", ErrorCouldNotAuthenticate)
}
