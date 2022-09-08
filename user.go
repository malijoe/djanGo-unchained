package django

type User struct {
	ID              uint
	Name, Password  string
	isAuthenticated bool
	isAdmin         bool
}

func (u *User) ValidatePassword(password string) error {
	return nil
}

func (u *User) IsAuthenticated() bool {
	return u.isAuthenticated
}

func (u *User) IsAdmin() bool {
	return u.isAdmin
}
