package domain

type UserLogin struct {
	Username string
	Password string
}

type TokenResponse struct {
	Token string `json:"token"`
}
