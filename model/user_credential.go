package model

type UserCredential struct {
	Username string `json:"userName"`
	Password string `json:"userPassword"`
	Email    string
}
