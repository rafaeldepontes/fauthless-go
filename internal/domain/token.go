package domain

import "time"

type TokenResponse struct {
	Token string `json:"token"`
}

type TokenRefreshResponse struct {
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	SessionId             string    `json:"session_id"`
	AccessToken           string    `json:"access_token"`
	RefreshToken          string    `json:"refresh_token"`
}

type RenewAccessTokenRequest struct {
	RefreshToken string `json:"access_token"`
}

type RenewAccessTokenResponse struct {
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
	AccessToken          string    `json:"access_token"`
}
