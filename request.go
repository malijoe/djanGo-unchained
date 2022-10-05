package django

import "net/http"

type Request struct {
	User           User
	Authenticators []AuthenticationClass
	Auth           string
	*http.Request
}

func (r *Request) Session() {

}
