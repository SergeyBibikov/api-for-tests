package main

type GetTokenBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ValidateTokenBody struct {
	Token string `json:"token"`
}
