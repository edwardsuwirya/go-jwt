package delivery

import (
	"enigmacamp.com/gojwt/config"
	"enigmacamp.com/gojwt/delivery/controller"
	"enigmacamp.com/gojwt/delivery/middleware"
	"enigmacamp.com/gojwt/usecase"
	"enigmacamp.com/gojwt/utils/authenticator"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type Server struct {
	engine       *gin.Engine
	host         string
	authUseCase  usecase.AuthUseCase
	tokenService authenticator.AccessToken
}

func (s *Server) Run() {
	s.initController()
	err := s.engine.Run(s.host)
	if err != nil {
		panic(err)
	}
}
func (s *Server) initController() {
	publicRoute := s.engine.Group("/enigma")
	tokenMdw := middleware.NewTokenValidator(s.tokenService)
	controller.NewAppController(publicRoute, s.authUseCase, tokenMdw)
}
func NewServer() *Server {
	c := config.NewConfig()
	r := gin.Default()
	client := redis.NewClient(&redis.Options{
		Addr:     c.RedisConfig.Address,
		Password: c.RedisConfig.Password,
		DB:       c.RedisConfig.Db,
	})
	tokenService := authenticator.NewAccessToken(c.TokenConfig, client)
	authUserCase := usecase.NewAuthUseCase(tokenService)
	if c.ApiHost == "" || c.ApiPort == "" {
		panic("No Host or port define")
	}
	host := fmt.Sprintf("%s:%s", c.ApiHost, c.ApiPort)
	return &Server{engine: r, host: host, authUseCase: authUserCase, tokenService: tokenService}
}
