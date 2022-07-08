package middleware

import (
	"enigmacamp.com/gojwt/authenticator"
	"fmt"
	"github.com/gin-gonic/gin"
)

type authCookieMiddleware struct {
	acctToken authenticator.Token
}

func NewTokenCookieValidator(acctToken authenticator.Token) AuthTokenMiddleware {
	return &authCookieMiddleware{
		acctToken: acctToken,
	}
}

func (a *authCookieMiddleware) RequireToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, _ := c.Cookie("token")
		fmt.Println("Cookie")
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
