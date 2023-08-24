package repository

import "english/cmd/domain/model"

type UserRepository interface {
	FindByEmail(email string) (*model.User, error)
	FindByIssAndSub(iss, sub string) (*model.User, error)
	Create(user *model.User) error
	Update(user *model.User) error
}
