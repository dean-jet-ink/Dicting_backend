package gateway

import (
	"english/cmd/domain/model"
	"english/cmd/domain/repository"
	"english/cmd/infrastructure/entity"
	"english/myerror"
	"errors"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type UserMySQLRepository struct {
	db *gorm.DB
}

func NewUserMySQLRepository(db *gorm.DB) repository.UserRepository {
	return &UserMySQLRepository{
		db: db,
	}
}

func (ur *UserMySQLRepository) FindByEmail(email string) (*model.User, error) {
	ue := &entity.UserEntity{}
	if err := ur.db.Where("email=?", email).First(ue).Error; err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return nil, myerror.ErrRecordNotFound
		}

		return nil, err
	}

	u := model.NewUser(ue.Id, ue.Email, ue.Password, ue.Name, ue.ProfileImageURL)

	return u, nil
}

func (ur *UserMySQLRepository) Create(user *model.User) error {
	entity, err := ur.modelToEntity(user)
	if err != nil {
		return err
	}

	if err := ur.db.Create(entity).Error; err != nil {
		return err
	}

	ur.entiryToModel(entity, user)

	return nil
}

func (ur *UserMySQLRepository) modelToEntity(u *model.User) (*entity.UserEntity, error) {
	e := &entity.UserEntity{}
	err := copier.Copy(e, u)
	if err != nil {
		return nil, err
	}

	return e, nil
}

func (ur *UserMySQLRepository) entiryToModel(e *entity.UserEntity, u *model.User) {
	u.SetId(e.Id)
	u.SetEmail(e.Email)
	u.SetPassword(e.Password)
	u.SetName(e.Name)
	u.SetProfileImageURL(e.ProfileImageURL)
}
