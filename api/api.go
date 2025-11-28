package api

import (
	"database/sql"
	"os"

	"github.com/joho/godotenv"
	"github.com/rafaeldepontes/auth-go/configs"
	"github.com/rafaeldepontes/auth-go/internal/database"
	"github.com/rafaeldepontes/auth-go/internal/repository"
	"github.com/rafaeldepontes/auth-go/internal/service"
	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

const (
	_ = iota
	CookieBased
	JwtBased
	JwtRefreshBased
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

// Init initialize all the resources needed for the server run properly
func Init() (*configs.Configuration, *Application, *sql.DB, error) {
	var logger *log.Logger = initLogger()
	godotenv.Load(".env", ".env.example")

	config := &configs.Configuration{
		JwtBasedPort:        os.Getenv("JWT_PORT"),
		CookieBasedPort:     os.Getenv("COOKIE_PORT"),
		JwtRefreshBasedPort: os.Getenv("JWT_REFRESH_PORT"),
	}

	db, err := database.Open()

	var userRepository *repository.UserRepository = repository.NewUserRepository(db)
	var userService *service.UserService = service.NewUserService(userRepository, logger)

	application := &Application{
		UserService: userService,
		Logger:      logger,
	}

	return config, application, db, err
}
