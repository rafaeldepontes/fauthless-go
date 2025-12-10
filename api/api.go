package api

import (
	"database/sql"
	"os"

	"github.com/joho/godotenv"

	"github.com/rafaeldepontes/fauthless-go/configs"
	"github.com/rafaeldepontes/fauthless-go/internal/auth"
	authRepository "github.com/rafaeldepontes/fauthless-go/internal/auth/repository"
	authServer "github.com/rafaeldepontes/fauthless-go/internal/auth/server"
	authService "github.com/rafaeldepontes/fauthless-go/internal/auth/service"
	"github.com/rafaeldepontes/fauthless-go/internal/cache"
	"github.com/rafaeldepontes/fauthless-go/internal/middleware"
	"github.com/rafaeldepontes/fauthless-go/internal/user"
	userRepository "github.com/rafaeldepontes/fauthless-go/internal/user/repository"
	userServer "github.com/rafaeldepontes/fauthless-go/internal/user/server"
	userService "github.com/rafaeldepontes/fauthless-go/internal/user/service"
	"github.com/rafaeldepontes/fauthless-go/pkg/db/postgres"

	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

const (
	_ = iota
	CookieBased
	JwtBased
	JwtRefreshBased
	OAuth2
)

func initLogger() *log.Logger {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	return &log.Logger{
		Out:   os.Stderr,
		Level: log.DebugLevel,
		Formatter: &prefixed.TextFormatter{
			DisableColors:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
			ForceFormatting: true,
		},
	}
}

// Init initialize all the resources needed for the server run properly.
func Init() (*configs.Configuration, *Application, *sql.DB, error) {
	var logger *log.Logger = initLogger()
	godotenv.Load(".env", ".env.example")

	config := &configs.Configuration{
		JwtSecretKey:        os.Getenv("JWT_SECRET_KEY"),
		JwtBasedPort:        os.Getenv("JWT_PORT"),
		CookieBasedPort:     os.Getenv("COOKIE_PORT"),
		OAuth2Port:          os.Getenv("OAUTH2_PORT"),
		JwtRefreshBasedPort: os.Getenv("JWT_REFRESH_PORT"),
		GoogleClientId:      os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleSecretKey:     os.Getenv("GOOGLE_KEY"),
		GoogleClientSecret:  os.Getenv("GOOGLE_CLIENT_SECRET"),
		UrlCallback:         os.Getenv("URL_CALLBACK"),
	}

	auth.InitOAuth(config)

	db, err := postgres.Open()

	var caches *cache.Caches = cache.NewCacheStorage()

	var userRepository user.Repository = userRepository.NewUserRepository(db)
	var sessionRepository auth.Repository = authRepository.NewSessionRepository(db)

	var userService user.Service = userService.NewUserService(userRepository, logger, caches)
	var authService auth.Service = authService.NewAuthService(userRepository, sessionRepository, logger, config.JwtSecretKey, caches)

	var userController user.Controller = userServer.NewUserController(&userService)
	var authController auth.Controller = authServer.NewAuthController(&authService)

	var middleware *middleware.Middleware = middleware.NewMiddleware(config.JwtSecretKey, caches)

	application := &Application{
		UserController: &userController,
		AuthController: &authController,
		Middleware:     middleware,
		Logger:         logger,
	}

	return config, application, db, err
}
