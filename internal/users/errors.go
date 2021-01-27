package users

// WrongUsernameOrPasswordError enable dependency injection for auth errors
type WrongUsernameOrPasswordError struct{}

func (m *WrongUsernameOrPasswordError) Error() string {
	return "wrong username or password"
}
