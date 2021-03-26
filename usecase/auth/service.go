package auth

import (
	"jwtauthv2/entity"
)

type AuthService struct {
	userStore  IUserStore
	tokenStore ITokenStore
}

func NewAuthService(r IUserStore, t ITokenStore) *AuthService {
	return &AuthService{
		userStore:  r,
		tokenStore: t,
	}
}

func (s *AuthService) SignUp(name, password string) (*entity.ID, error) {

	user, err := entity.NewUser(name, password)

	if err != nil {
		return nil, err
	}

	if err = s.userStore.Create(user); err != nil {
		return nil, err
	}

	return &user.ID, nil
}

func (s *AuthService) SignIn(name, password string) (*entity.TokenPair, error) {

	user, err := s.userStore.FindByName(name)

	if err != nil {
		return nil, err
	}

	if err = user.ComparePassword(password); err != nil {
		return nil, err
	}

	tokenPair, err := s.issueTokenPair(&user.ID)

	if err != nil {
		return nil, err
	}

	return tokenPair, nil
}

func (s *AuthService) Refresh(token string) (*entity.TokenPair, error) {

	user, err := s.tokenStore.FindByToken(token)

	if err != nil {
		return nil, err
	}

	err = s.tokenStore.RemoveByToken(token)

	if err != nil {
		return nil, err
	}

	tokenPair, err := s.issueTokenPair(user)

	if err != nil {
		return nil, err
	}

	return tokenPair, nil
}

func (s *AuthService) Logout(user *entity.ID) error {
	return s.tokenStore.RemoveByID(user)
}

func (s *AuthService) issueTokenPair(user *entity.ID) (*entity.TokenPair, error) {
	a, err := entity.IssueToken(user, entity.AccessTokenTTL)

	if err != nil {
		return nil, err
	}

	r, err := entity.IssueToken(user, entity.RefreshTokenTTL)

	if err != nil {
		return nil, err
	}

	err = s.tokenStore.Add(user, r)

	if err != nil {
		return nil, err
	}

	return &entity.TokenPair{Access: a, Refresh: r}, nil
}
