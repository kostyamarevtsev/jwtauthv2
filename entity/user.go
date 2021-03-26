package entity

import (
	"crypto/sha1"
	"fmt"
	"jwtauthv2"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
)

type ID = uuid.UUID

type User struct {
	ID       ID     `db:"id" json:"id"`
	Name     string `db:"name" json:"name"`
	Password string `db:"password" json:"-"`
}

func NewID() ID {
	return uuid.New()
}

func StringToID(s string) (ID, error) {
	return uuid.Parse(s)
}

func NewUser(name, password string) (*User, error) {
	u := &User{
		ID: NewID(),
	}

	if err := u.SetName(name); err != nil {
		return u, err
	}

	if err := u.SetPassword(password); err != nil {
		return u, err
	}

	return u, nil
}

func (u *User) SetName(s string) error {

	err := validation.Validate(s,
		validation.Required,
		validation.Length(3, 30))

	if err != nil {
		return &jwtauthv2.Error{
			Code:    jwtauthv2.EINVALID,
			Message: "Invalid name",
		}
	}

	u.Name = s

	return nil
}

func (u *User) SetPassword(s string) error {
	err := validation.Validate(s,
		validation.Required,
		validation.Length(5, 40))

	if err != nil {
		return &jwtauthv2.Error{
			Code:    jwtauthv2.EINVALID,
			Message: "Invalid password",
		}
	}

	u.Password = GetHashString(s)

	return nil
}

func (u *User) ComparePassword(s string) error {

	if u.Password != GetHashString(s) {
		return &jwtauthv2.Error{
			Code:    jwtauthv2.EINVALID,
			Message: "Password is wrong",
		}
	}

	return nil
}

func GetHashString(s string) string {
	hash := sha1.New()
	hash.Write([]byte(s))
	return fmt.Sprintf("%x", hash.Sum(nil))
}
