package envirotment

import "os"

type AppEnv string

const (
	Dev AppEnv = "development"
	Stg AppEnv = "staging"
	Prd AppEnv = "production"
)

var app_env_aliases map[string]AppEnv = map[string]AppEnv{
	"dev":         Dev,
	"stg":         Stg,
	"prd":         Prd,
	"prod":        Prd,
	"development": Dev,
	"staging":     Stg,
	"production":  Prd,
}

const (
	// ENV
	APP_ENV = "APP_ENV"

	// APP
	SECRET_KEY        = "SECRET_KEY"
	MAIL_PORT         = "MAIL_PORT"
	MAIL_HOST         = "MAIL_HOST"
	MAIL_FROM_ADDRESS = "MAIL_FROM_ADDRESS"
	MAIL_PASSWORD     = "MAIL_PASSWORD"

	// DATABASE

	DATABASE_HOST     = "DATABASE_HOST"
	DATABASE_USER     = "DATABASE_USER"
	DATABASE_PASSWORD = "DATABASE_PASSWORD"
	DATABASE_PORT     = "DATABASE_PORT"
	DATABASE_DATABASE = "DATABASE_DATABASE"
	DATABASE_URL      = "DATABASE_URL"
)

func Get(key string) string {
	return os.Getenv(key)
}

func GetSecretKey() string {
	return os.Getenv(SECRET_KEY)
}

func GetAppEnv() AppEnv {

	app_env := os.Getenv(APP_ENV)

	if app_env_alias, ok := app_env_aliases[app_env]; ok {
		return app_env_alias
	}

	return Dev
}
