package services

type UserService interface {
	Signup(email, username, password string) error
}
