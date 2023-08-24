package repository

import "english/cmd/domain/model"

type UserRepository interface {
	FindById(id string) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	FindByIssAndSub(iss, sub string) (*model.User, error)
	Create(user *model.User) error
	Update(user *model.User) error
}
