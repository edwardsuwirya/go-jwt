package config

import (
	"github.com/golang-jwt/jwt"
	"time"
)

type ApiConfig struct {
	ApiPort string
	ApiHost string
}

type TokenConfig struct {
	ApplicationName     string
	JwtSignatureKey     string
	JwtSigningMethod    *jwt.SigningMethodHMAC
	AccessTokenLifeTime time.Duration
}

type Config struct {
	ApiConfig
	TokenConfig
}

func (c Config) readConfig() Config {
	c.ApiConfig = ApiConfig{
		ApiPort: "8888",
		ApiHost: "localhost",
	}
	c.TokenConfig = TokenConfig{
		ApplicationName:     "ENIGMA",
		JwtSignatureKey:     "P@ssw0rd",
		JwtSigningMethod:    jwt.SigningMethodHS256,
		AccessTokenLifeTime: 60 * time.Second,
	}
	return c
}
func NewConfig() Config {
	cfg := Config{}
	return cfg.readConfig()
}
