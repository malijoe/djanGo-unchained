package django

type User struct {
	ID             uint
	Name, Password string
}

func (u *User) ValidatePassword(password string) error {
	return nil
}
