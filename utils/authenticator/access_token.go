package authenticator

import (
	"context"
	"enigmacamp.com/gojwt/config"
	"enigmacamp.com/gojwt/model"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"log"
	"time"
)

type AccessToken interface {
	CreateAccessToken(cred *model.UserCredential) (TokenDetail, error)
	VerifyAccessToken(tokenString string) (AccessDetails, error)
	StoreAccessToken(userName string, tokenDetail TokenDetail) error
	FetchAccessToken(accessDetails AccessDetails) error
	DeleteAccessToken(accessUUID string) error
}

type accessToken struct {
	cfg    config.TokenConfig
	client *redis.Client
}

func NewAccessToken(config config.TokenConfig, client *redis.Client) AccessToken {
	return &accessToken{
		cfg:    config,
		client: client,
	}
}
func (t *accessToken) StoreAccessToken(userName string, tokenDetail TokenDetail) error {
	at := time.Unix(tokenDetail.AtExpires, 0)
	now := time.Now()
	err := t.client.Set(context.Background(), tokenDetail.AccessUuid, userName, at.Sub(now)).Err()
	if err != nil {
		return err
	}
	return nil
}

func (t *accessToken) FetchAccessToken(accessDetails AccessDetails) error {
	userName, err := t.client.Get(context.Background(), accessDetails.AccessUuid).Result()
	if err != nil {
		return err
	}
	if userName == "" {
		return errors.New("Invalid token")
	}
	return nil
}

func (t *accessToken) DeleteAccessToken(accessUUID string) error {
	rowAffected, err := t.client.Del(context.Background(), accessUUID).Result()
	if err != nil {
		return err
	}
	if rowAffected == 0 {
		return errors.New("Invalid token")
	}
	return nil
}

func (t *accessToken) CreateAccessToken(cred *model.UserCredential) (TokenDetail, error) {
	td := TokenDetail{}
	now := time.Now().UTC()
	end := now.Add(t.cfg.AccessTokenLifeTime)
	td.AtExpires = end.Unix()
	td.AccessUuid = uuid.New().String()

	claims := MyClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer: t.cfg.ApplicationName,
		},
		Username:   cred.Username,
		Email:      cred.Email,
		AccessUUID: td.AccessUuid,
	}
	claims.IssuedAt = now.Unix()
	claims.ExpiresAt = end.Unix()
	token := jwt.NewWithClaims(
		t.cfg.JwtSigningMethod,
		claims,
	)
	newToken, err := token.SignedString([]byte(t.cfg.JwtSignatureKey))
	if err != nil {
		return td, err
	}
	td.AccessToken = newToken
	return td, nil
}

func (t *accessToken) VerifyAccessToken(tokenString string) (AccessDetails, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Signing method invalid")
		} else if method != t.cfg.JwtSigningMethod {
			return nil, fmt.Errorf("Signing method invalid")
		}

		return []byte(t.cfg.JwtSignatureKey), nil
	})
	claims, ok := token.Claims.(jwt.MapClaims)
	accessDetails := AccessDetails{}
	if !ok || !token.Valid || claims["iss"] != t.cfg.ApplicationName {
		log.Println("Token Invalid")
		return accessDetails, err
	}
	userName := claims["Username"].(string)
	accessDetails.AccessUuid = claims["AccessUUID"].(string)
	accessDetails.UserName = userName
	return accessDetails, nil
}
