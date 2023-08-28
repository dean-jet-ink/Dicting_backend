package gateway

import (
	"english/cmd/domain/model"
	"english/cmd/domain/repository"
	"english/cmd/infrastructure/entity"
	"english/myerror"
	"errors"
	"strings"

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

func (ur *UserMySQLRepository) FindById(id string) (*model.User, error) {
	entity := &entity.UserEntity{}
	if err := ur.db.Where("id=?", id).First(entity).Error; err != nil {
		return nil, err
	}

	user := model.NewUser(entity.Id, entity.Email, entity.Password, entity.Name, entity.ProfileImageURL)

	return user, nil
}

func (ur *UserMySQLRepository) FindByEmail(email string) (*model.User, error) {
	entity := &entity.UserEntity{}
	if err := ur.db.Where("email=?", email).First(entity).Error; err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return nil, myerror.ErrRecordNotFound
		}

		return nil, err
	}

	u := model.NewUser(entity.Id, entity.Email, entity.Password, entity.Name, entity.ProfileImageURL)

	return u, nil
}

func (ur *UserMySQLRepository) FindByIssAndSub(iss, sub string) (*model.User, error) {
	entity := &entity.UserEntity{}
	if err := ur.db.Where("iss=? AND sub=?", iss, sub).First(entity).Error; err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return nil, myerror.ErrRecordNotFound
		}

		return nil, err
	}

	user := model.NewUser(entity.Id, entity.Email, entity.Password, entity.Name, entity.ProfileImageURL)

	return user, nil
}

func (ur *UserMySQLRepository) Create(user *model.User) error {
	entity, err := ur.modelToEntity(user)
	if err != nil {
		return err
	}

	if err := ur.db.Create(entity).Error; err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			return myerror.ErrDuplicatedKey
		}
		return err
	}

	ur.entiryToModel(entity, user)

	return nil
}

func (ur *UserMySQLRepository) Update(user *model.User) error {
	entity := &entity.UserEntity{
		Id:              user.Id(),
		Email:           user.Email(),
		Name:            user.Name(),
		ProfileImageURL: user.ProfileImageURL(),
	}

	if err := ur.db.Updates(entity).Error; err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			return myerror.ErrDuplicatedKey
		}
		return err
	}

	return nil
}

func (ur *UserMySQLRepository) modelToEntity(m *model.User) (*entity.UserEntity, error) {
	e := &entity.UserEntity{}
	err := copier.Copy(e, m)
	if err != nil {
		return nil, err
	}

	return e, nil
}

func (ur *UserMySQLRepository) entiryToModel(e *entity.UserEntity, m *model.User) {
	m.SetId(e.Id)
	m.SetEmail(e.Email)
	m.SetPassword(e.Password)
	m.SetName(e.Name)
	m.SetProfileImageURL(e.ProfileImageURL)
}
