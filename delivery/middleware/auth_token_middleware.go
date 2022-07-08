package middleware

import (
	"enigmacamp.com/gojwt/authenticator"
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
	acctToken authenticator.Token
}

func NewTokenValidator(acctToken authenticator.Token) AuthTokenMiddleware {
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
		token, err := a.acctToken.VerifyAccessToken(tokenString)
		if err != nil {
			c.JSON(401, gin.H{
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}
		fmt.Println(token)
		if token != nil {
			c.Next()
		} else {
			c.JSON(401, gin.H{
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}
	}
}
