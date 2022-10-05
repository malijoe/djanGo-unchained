package django

type User interface {
	ValidatePassword(password string) error
	IsAuthenticated() bool
	IsAdmin() bool
	Username() string
}
