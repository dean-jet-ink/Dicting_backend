package model

import (
	"english/config"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type User struct {
	id            string
	email         string
	password      string
	name          string
	profileImgURL string
	iss           string
	sub           string
}

func NewUser(id, email, password, name, profileImageURL string) *User {
	return &User{
		id:            id,
		email:         email,
		password:      password,
		name:          name,
		profileImgURL: profileImageURL,
	}
}

func (u *User) CreateJWT(expireSec int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": u.id,
		"exp":     time.Now().Add(time.Second * time.Duration(expireSec)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 環境変数のSECRETを使用し署名
	jwtToken, err := token.SignedString([]byte(config.JWTSecret()))
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
	return u.profileImgURL
}

func (u *User) Iss() string {
	return u.iss
}

func (u *User) Sub() string {
	return u.sub
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
	u.profileImgURL = profileImageURL
}

func (u *User) SetIss(iss string) {
	u.iss = iss
}

func (u *User) SetSub(sub string) {
	u.sub = sub
}
