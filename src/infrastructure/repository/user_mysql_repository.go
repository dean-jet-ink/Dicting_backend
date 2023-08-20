package repository

import (
	"english/myerror"
	"english/src/domain/usermodel"
	"english/src/infrastructure/entity"
	"errors"

	"gorm.io/gorm"
)

type UserMySQLRepository struct {
	db *gorm.DB
}

func NewUserMySQLRepository(db *gorm.DB) usermodel.UserRepository {
	return &UserMySQLRepository{
		db: db,
	}
}

func (ur *UserMySQLRepository) FindByEmail(email string) (*usermodel.User, error) {
	ue := &entity.UserEntity{}
	if err := ur.db.Where("email=?", email).First(ue).Error; err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return nil, myerror.ErrRecordNotFound
		}

		return nil, err
	}

	u := usermodel.NewUser(ue.Id, ue.Email, ue.Password, ue.Name, ue.ProfileImageURL)

	return u, nil
}
