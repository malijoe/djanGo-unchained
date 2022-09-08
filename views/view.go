package views

import (
	"github.com/malijoe/djanGo-unchained/authentication"
	"github.com/malijoe/djanGo-unchained/permissions"
)

type View struct {
	AuthenticationClasses []authentication.Class
	PermissionClasses     []permissions.Class
	AllowedMethods        []string
}
