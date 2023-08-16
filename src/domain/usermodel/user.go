package usermodel

type User struct {
	id              int
	email           string
	password        string
	name            string
	profileImageURL string
}

func NewUser(id int, email, password, name, profileImageURL string) *User {
	return &User{
		id:              id,
		email:           email,
		password:        password,
		name:            name,
		profileImageURL: profileImageURL,
	}
}

func (u *User) Id() int {
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
