package api

import (
	"database/sql"
	"os"

	"github.com/joho/godotenv"

	"github.com/rafaeldepontes/auth-go/configs"
	"github.com/rafaeldepontes/auth-go/internal/auth"
	"github.com/rafaeldepontes/auth-go/internal/database"
	"github.com/rafaeldepontes/auth-go/internal/database/repository"
	"github.com/rafaeldepontes/auth-go/internal/middleware"
	"github.com/rafaeldepontes/auth-go/internal/service"
	"github.com/rafaeldepontes/auth-go/internal/storage"

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

	db, err := database.Open()

	var caches *storage.Caches = storage.NewCacheStorage()

	var userRepository *repository.UserRepository = repository.NewUserRepository(db)
	var sessionRepository *repository.SessionRepository = repository.NewSessionRepository(db)

	var userService *service.UserService = service.NewUserService(userRepository, logger, caches)
	var authService *service.AuthService = service.NewAuthService(userRepository, sessionRepository, logger, config.JwtSecretKey, caches)

	var middleware *middleware.Middleware = middleware.NewMiddleware(config.JwtSecretKey, caches)

	application := &Application{
		UserService: userService,
		AuthService: authService,
		Middleware:  middleware,
		Logger:      logger,
	}

	return config, application, db, err
}
