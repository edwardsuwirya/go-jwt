package authenticator

import "github.com/golang-jwt/jwt"

type MyClaims struct {
	jwt.StandardClaims
	Username   string `json:"Username"`
	Email      string `json:"Email"`
	AccessUUID string
}
