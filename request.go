package django

import "net/http"

type Request struct {
	User *User
	Auth string
	*http.Request
	http.ResponseWriter
}

func (r *Request) Session() {

}
