package middleware

import (
	"enigmacamp.com/gojwt/utils/authenticator"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

type authHeader struct {
	AuthorizationHeader string `header:"Authorization"`
}
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
		h := authHeader{}
		if err := c.ShouldBindHeader(&h); err != nil {
			c.JSON(401, gin.H{
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}
		tokenString := strings.Replace(h.AuthorizationHeader, "Bearer ", "", -1)
		fmt.Println(tokenString)
		if tokenString == "" {
			c.JSON(401, gin.H{
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}
		accountDetail, err := a.acctToken.VerifyAccessToken(tokenString)
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
