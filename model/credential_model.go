package model

type CredentialModel struct {
	Username string `json:"userName"`
	Password string `json:"userPassword"`
	Email    string
}
