package controller

import (
	"enigmacamp.com/gojwt/delivery/middleware"
	"enigmacamp.com/gojwt/model"
	"enigmacamp.com/gojwt/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AppController struct {
	rg          *gin.RouterGroup
	authUseCase usecase.AuthUseCase
}

func (cc *AppController) getAllCustomer(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "user",
	})
}
func (cc *AppController) userAuth(ctx *gin.Context) {
	var user model.UserCredential
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "can't bind struct",
		})
		return
	}
	token, err := cc.authUseCase.UserAuth(user)
	if err != nil {
		ctx.AbortWithStatus(401)
		return
	}
	ctx.JSON(200, gin.H{
		"token": token,
	})
}

func NewAppController(routerGroup *gin.RouterGroup, authUseCase usecase.AuthUseCase, tokenMdw middleware.AuthTokenMiddleware) *AppController {
	ctrl := AppController{
		rg:          routerGroup,
		authUseCase: authUseCase,
	}
	ctrl.rg.POST("/auth", ctrl.userAuth)
	protectedGroup := ctrl.rg.Group("/protected", tokenMdw.RequireToken())
	protectedGroup.GET("/user", ctrl.getAllCustomer)
	return &ctrl
}
