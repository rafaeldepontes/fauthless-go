package api

import (
	"os"

	"github.com/joho/godotenv"
)

func Init() *Configuration {
	godotenv.Load(".env", ".env.example")

	config := &Configuration{
		JwtBasedPort:        os.Getenv("JWT_PORT"),
		CookieBasedPort:     os.Getenv("COOKIE_PORT"),
		JwtRefreshBasedPort: os.Getenv("JWT_REFRESH_PORT"),
	}

	return config
}
