package authentication

import "errors"

var (
	ErrorCouldNotAuthenticate = errors.New("could not authenticate")
)

type AuthenticationError struct {
	status int
}
