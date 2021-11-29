package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type authHeader struct {
	AuthorizationHeader string `header:"Authorization"`
}

type Credential struct {
	Username string
	Password string
}

func main() {
	r := gin.Default()
	r.Use(AuthTokenMiddleware())
	r.POST("/login", func(c *gin.Context) {
		var user Credential
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "can't bind struct",
			})
			return
		}
		if user.Username == "enigma" && user.Password == "123" {
			c.JSON(200, gin.H{
				"token": "123",
			})
		} else {
			c.AbortWithStatus(401)
		}

	})
	r.GET("/customer", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "customer",
		})
	})
	r.GET("/product", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "product",
		})
	})
	err := r.Run("localhost:8888")
	if err != nil {
		panic(err)
	}
}
func AuthTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/login" {
			c.Next()
			fmt.Println("sss")
		} else {
			h := authHeader{}
			if err := c.ShouldBindHeader(&h); err != nil {
				c.JSON(401, gin.H{
					"message": "Unauthorized",
				})
				c.Abort()
			}
			if h.AuthorizationHeader == "123" {
				c.Next()
			} else {
				c.JSON(401, gin.H{
					"message": "Unauthorized",
				})
				c.Abort()
			}
		}
	}
}
