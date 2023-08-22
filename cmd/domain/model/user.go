package model

import (
	"english/config"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type User struct {
	id              string
	email           string
	password        string
	name            string
	profileImageURL string
}

func NewUser(id, email, password, name, profileImageURL string) *User {
	return &User{
		id:              id,
		email:           email,
		password:        password,
		name:            name,
		profileImageURL: profileImageURL,
	}
}

func (u *User) CreateJWT(expireSec int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": u.id,
		"exp":     time.Now().Add(time.Second * time.Duration(expireSec)).Unix(),
	})
	// 環境変数のSECRETを使用し署名
	jwtToken, err := token.SignedString([]byte(config.Secret()))
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}

func (u *User) Id() string {
	return u.id
}

func (u *User) Email() string {
	return u.email
}

func (u *User) Password() string {
	return u.password
}

func (u *User) Name() string {
	return u.name
}

func (u *User) ProfileImageURL() string {
	return u.profileImageURL
}

func (u *User) SetId(id string) {
	u.id = id
}

func (u *User) SetEmail(email string) {
	u.email = email
}

func (u *User) SetPassword(password string) {
	u.password = password
}

func (u *User) SetName(name string) {
	u.name = name
}

func (u *User) SetProfileImageURL(profileImageURL string) {
	u.profileImageURL = profileImageURL
}
