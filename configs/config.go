package configs

type Configuration struct {
	SecretKey           string
	JwtBasedPort        string
	CookieBasedPort     string
	JwtRefreshBasedPort string
}
