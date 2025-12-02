package configs

type Configuration struct {
	JwtSecretKey        string
	JwtBasedPort        string
	CookieBasedPort     string
	JwtRefreshBasedPort string
	OAuth2Port          string
	GoogleClientId      string
	GoogleSecretKey     string
	GoogleClientSecret  string
	UrlCallback         string
}
