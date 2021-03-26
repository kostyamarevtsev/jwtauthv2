package auth

import "jwtauthv2/entity"

type IUserStore interface {
	Create(user *entity.User) error
	FindByName(name string) (*entity.User, error)
}

type ITokenStore interface {
	Add(id *entity.ID, token string) error
	FindByToken(token string) (*entity.ID, error)
	RemoveByToken(token string) error
	RemoveByID(id *entity.ID) error
}

type IAuthService interface {
	SignUp(name, password string) (*entity.ID, error)
	SignIn(name, password string) (*entity.TokenPair, error)
	Refresh(token string) (*entity.TokenPair, error)
	Logout(user *entity.ID) error
}
