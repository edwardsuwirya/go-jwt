package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	type authHeader struct {
		AuthorizationHeader string `header:"Authorization"`
	}

	r.GET("/customer", func(c *gin.Context) {
		h := authHeader{}
		if err := c.ShouldBindHeader(&h); err != nil {
			c.JSON(401, gin.H{
				"message": "Unauthorized",
			})
			return
		}
		if h.AuthorizationHeader == "123" {
			c.JSON(200, gin.H{
				"message": "customer",
			})
			return
		}
		c.JSON(401, gin.H{
			"message": "Unauthorized",
		})
	})
	r.GET("/product", func(c *gin.Context) {
		h := authHeader{}
		if err := c.ShouldBindHeader(&h); err != nil {
			c.JSON(401, gin.H{
				"message": "Unauthorized",
			})
			return
		}
		if h.AuthorizationHeader == "123" {
			c.JSON(200, gin.H{
				"message": "product",
			})
			return
		}
		c.JSON(401, gin.H{
			"message": "Unauthorized",
		})
	})
	err := r.Run("localhost:8888")
	if err != nil {
		panic(err)
	}
}
