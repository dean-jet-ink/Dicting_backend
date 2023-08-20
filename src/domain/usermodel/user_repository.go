package usermodel

type UserRepository interface {
	FindByEmail(email string) (*User, error)
}
