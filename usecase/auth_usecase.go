package usecase

import (
	"enigmacamp.com/gojwt/model"
	"enigmacamp.com/gojwt/utils/authenticator"
)

type AuthUseCase interface {
	UserAuth(user model.UserCredential) (token string, err error)
}

type authUseCase struct {
	tokenService authenticator.AccessToken
}

func (a *authUseCase) UserAuth(user model.UserCredential) (token string, err error) {
	if user.Username == "enigma" && user.Password == "123" {
		token, err := a.tokenService.CreateAccessToken(&user)
		if err != nil {
			return "", err
		}
		return token, nil
	} else {
		return "", err
	}
}

func NewAuthUseCase(service authenticator.AccessToken) AuthUseCase {
	authUseCase := new(authUseCase)
	authUseCase.tokenService = service
	return authUseCase
}
