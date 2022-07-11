package authenticator

import (
	"enigmacamp.com/gojwt/config"
	"enigmacamp.com/gojwt/model"
	"fmt"
	"github.com/golang-jwt/jwt"
	"log"
	"time"
)

type Token interface {
	CreateAccessToken(cred *model.UserCredential) (string, error)
	VerifyAccessToken(tokenString string) (jwt.MapClaims, error)
}

type token struct {
	cfg config.TokenConfig
}

func NewTokenService(config config.TokenConfig) Token {
	return &token{
		cfg: config,
	}
}

func (t *token) CreateAccessToken(cred *model.UserCredential) (string, error) {
	now := time.Now().UTC()
	end := now.Add(t.cfg.AccessTokenLifeTime)
	claims := MyClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer: t.cfg.ApplicationName,
		},
		Username: cred.Username,
		Email:    cred.Email,
	}
	claims.IssuedAt = now.Unix()
	claims.ExpiresAt = end.Unix()
	token := jwt.NewWithClaims(
		t.cfg.JwtSigningMethod,
		claims,
	)
	return token.SignedString([]byte(t.cfg.JwtSignatureKey))
}

func (t *token) VerifyAccessToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Signing method invalid")
		} else if method != t.cfg.JwtSigningMethod {
			return nil, fmt.Errorf("Signing method invalid")
		}

		return []byte(t.cfg.JwtSignatureKey), nil
	})
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid || claims["iss"] != t.cfg.ApplicationName {
		log.Println("Token Invalid")
		return nil, err
	}
	return claims, nil
}
