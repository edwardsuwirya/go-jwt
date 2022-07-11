package middleware

import (
	"enigmacamp.com/gojwt/utils/authenticator"
	"github.com/gin-gonic/gin"
)

type AuthTokenMiddleware interface {
	RequireToken() gin.HandlerFunc
}
type authTokenMiddleware struct {
	acctToken authenticator.AccessToken
}

func NewTokenValidator(acctToken authenticator.AccessToken) AuthTokenMiddleware {
	return &authTokenMiddleware{
		acctToken: acctToken,
	}
}

func (a *authTokenMiddleware) RequireToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := authenticator.BindAuthHeader(c)
		if err != nil {
			return
		}
		accountDetail, err := a.acctToken.VerifyAccessToken(token)
		if err != nil {
			c.JSON(401, gin.H{
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}

		// User sudah logout tetapi mengirimkan token yang sama
		err = a.acctToken.FetchAccessToken(accountDetail)
		if err != nil {
			c.JSON(401, gin.H{
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
