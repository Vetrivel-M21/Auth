package config

type Config struct {
	DB_USER    string
	DB_PASS    string
	DB_NAME    string
	DB_HOST    string
	DB_PORT    string
	SECRET_KEY string
}

var SecretKey string
