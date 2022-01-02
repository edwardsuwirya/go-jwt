package main

import (
	"enigmacamp.com/gojwt/authenticator"
	mdw "enigmacamp.com/gojwt/delivery/middleware"
	"enigmacamp.com/gojwt/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"time"
)

func main() {
	r := gin.Default()
	tokenConfig := authenticator.TokenConfig{
		ApplicationName:     "ENIGMA",
		JwtSignatureKey:     "P@ssw0rd",
		JwtSigningMethod:    jwt.SigningMethodHS256,
		AccessTokenLifeTime: 60 * time.Second,
	}
	tokenService := authenticator.NewTokenService(tokenConfig)
	r.Use(mdw.NewTokenValidator(tokenService).RequireToken())
	publicRoute := r.Group("/enigma")
	publicRoute.POST("/auth", func(c *gin.Context) {
		var user model.CredentialModel
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "can't bind struct",
			})
			return
		}
		if user.Username == "enigma" && user.Password == "123" {
			token, err := tokenService.CreateAccessToken(&user)
			if err != nil {
				c.AbortWithStatus(401)
			}
			c.JSON(200, gin.H{
				"token": token,
			})
		} else {
			c.AbortWithStatus(401)
		}

	})
	publicRoute.GET("/user", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "user",
		})
	})
	err := r.Run("localhost:8888")
	if err != nil {
		panic(err)
	}
}
