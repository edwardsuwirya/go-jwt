package main

import "github.com/gin-gonic/gin"

func main(){
	r := gin.Default()
	r.GET("/customer", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "customer",
		})
	})
	err := r.Run("localhost:8888")
	if err != nil {
		panic(err)
	}
}
