package authenticator

import (
	"enigmacamp.com/gojwt/delivery/request"
	"errors"
	"github.com/gin-gonic/gin"
	"strings"
)

func BindAuthHeader(c *gin.Context) (string, error) {
	h := new(request.AuthHeader)
	if err := c.ShouldBindHeader(h); err != nil {
		c.JSON(401, gin.H{
			"message": "Unauthorized",
		})
		c.Abort()
		return "", err
	}
	tokenString := strings.Replace(h.AuthorizationHeader, "Bearer ", "", -1)
	if tokenString == "Bearer" || tokenString == "" {
		c.JSON(401, gin.H{
			"message": "Unauthorized",
		})
		c.Abort()
		return "", errors.New("Empty token")
	}
	return tokenString, nil
}
