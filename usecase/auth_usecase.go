package usecase

import (
	"enigmacamp.com/gojwt/model"
	"enigmacamp.com/gojwt/utils/authenticator"
	"errors"
)

type AuthUseCase interface {
	UserAuth(user model.UserCredential) (token string, err error)
	Logout(token string) error
}

type authUseCase struct {
	tokenService authenticator.AccessToken
}

func (a *authUseCase) Logout(token string) error {
	accountDetail, err := a.tokenService.VerifyAccessToken(token)
	err = a.tokenService.DeleteAccessToken(accountDetail.AccessUuid)
	if err != nil {
		return errors.New("Token Invalid")
	}
	return nil
}

func (a *authUseCase) UserAuth(user model.UserCredential) (token string, err error) {
	if user.Username == "enigma" && user.Password == "123" {
		tokenDetail, err := a.tokenService.CreateAccessToken(&user)
		err = a.tokenService.StoreAccessToken(user.Username, tokenDetail)
		if err != nil {
			return "", err
		}
		return tokenDetail.AccessToken, nil
	} else {
		return "", err
	}
}

func NewAuthUseCase(service authenticator.AccessToken) AuthUseCase {
	authUseCase := new(authUseCase)
	authUseCase.tokenService = service
	return authUseCase
}
