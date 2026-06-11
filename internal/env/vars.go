package env

import (
	"fmt"
	"devdock/internal/config"
	"devdock/internal/utils"
)

func GetServiceVars(projName string, projType string, serviceName string, cfg *config.Config) map[string]string {
	vars := make(map[string]string)
	dbName := utils.NormalizeDBName(projName)

	if projType == "laravel" {
		if serviceName == "mysql" && cfg.Services.MySQL != nil {
			vars["DB_CONNECTION"] = "mysql"
			vars["DB_HOST"] = "mysql"
			vars["DB_PORT"] = fmt.Sprintf("%d", cfg.Services.MySQL.Port)
			vars["DB_DATABASE"] = dbName
			vars["DB_USERNAME"] = "devdock"
			vars["DB_PASSWORD"] = "password"
		} else if serviceName == "postgres" && cfg.Services.Postgres != nil {
			vars["DB_CONNECTION"] = "pgsql"
			vars["DB_HOST"] = "postgres"
			vars["DB_PORT"] = fmt.Sprintf("%d", cfg.Services.Postgres.Port)
			vars["DB_DATABASE"] = dbName
			vars["DB_USERNAME"] = "devdock"
			vars["DB_PASSWORD"] = "password"
		} else if serviceName == "redis" && cfg.Services.Redis != nil {
			vars["REDIS_HOST"] = "redis"
			vars["REDIS_PASSWORD"] = "null"
			vars["REDIS_PORT"] = fmt.Sprintf("%d", cfg.Services.Redis.Port)
		} else if serviceName == "mailpit" && cfg.Services.Mailpit != nil {
			vars["MAIL_MAILER"] = "smtp"
			vars["MAIL_HOST"] = "mailpit"
			vars["MAIL_PORT"] = fmt.Sprintf("%d", cfg.Services.Mailpit.SMTPPort)
		} else if serviceName == "minio" && cfg.Services.MinIO != nil {
			vars["AWS_ACCESS_KEY_ID"] = "devdock"
			vars["AWS_SECRET_ACCESS_KEY"] = "password"
			vars["AWS_DEFAULT_REGION"] = "us-east-1"
			vars["AWS_BUCKET"] = "local"
			vars["AWS_USE_PATH_STYLE_ENDPOINT"] = "true"
			vars["AWS_ENDPOINT"] = fmt.Sprintf("http://localhost:%d", cfg.Services.MinIO.APIPort)
		}
	} else if projType == "nextjs" || projType == "express" || projType == "fiber" {
		if serviceName == "postgres" && cfg.Services.Postgres != nil {
			vars["DATABASE_URL"] = fmt.Sprintf("postgresql://devdock:password@postgres:%d/%s", cfg.Services.Postgres.Port, dbName)
		} else if serviceName == "mysql" && cfg.Services.MySQL != nil {
			vars["DATABASE_URL"] = fmt.Sprintf("mysql://devdock:password@mysql:%d/%s", cfg.Services.MySQL.Port, dbName)
		} else if serviceName == "redis" && cfg.Services.Redis != nil {
			vars["REDIS_URL"] = fmt.Sprintf("redis://redis:%d", cfg.Services.Redis.Port)
		}
	}

	return vars
}
