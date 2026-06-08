package services

import (
	"fmt"
	"devdock/internal/config"
	"devdock/internal/utils"
)

type Service struct {
	Image       string            `yaml:"image"`
	Ports       []string          `yaml:"ports,omitempty"`
	Environment map[string]string `yaml:"environment,omitempty"`
	Volumes     []string          `yaml:"volumes,omitempty"`
	Networks    []string          `yaml:"networks,omitempty"`
	Healthcheck *Healthcheck      `yaml:"healthcheck,omitempty"`
	DependsOn   map[string]any    `yaml:"depends_on,omitempty"`
}

type Healthcheck struct {
	Test        []string `yaml:"test"`
	Interval    string   `yaml:"interval,omitempty"`
	Timeout     string   `yaml:"timeout,omitempty"`
	Retries     int      `yaml:"retries,omitempty"`
	StartPeriod string   `yaml:"start_period,omitempty"`
}

func BuildPostgres(projName string, cfg config.PostgresConfig) Service {
	dbName := utils.NormalizeDBName(projName)
	return Service{
		Image: fmt.Sprintf("postgres:%s", cfg.Version),
		Ports: []string{fmt.Sprintf("%d:5432", cfg.Port)},
		Environment: map[string]string{
			"POSTGRES_DB":       dbName,
			"POSTGRES_USER":     "devdock",
			"POSTGRES_PASSWORD": "password",
		},
		Volumes: []string{
			fmt.Sprintf("devdock_%s_postgres_data:/var/lib/postgresql/data", utils.NormalizeProjectName(projName)),
		},
		Networks: []string{"devdock_network"},
		Healthcheck: &Healthcheck{
			Test: []string{"CMD-SHELL", "pg_isready -U devdock -d " + dbName},
			Interval: "5s",
			Timeout: "5s",
			Retries: 5,
		},
	}
}

func BuildMySQL(projName string, cfg config.MySQLConfig) Service {
	dbName := utils.NormalizeDBName(projName)
	return Service{
		Image: fmt.Sprintf("mysql:%s", cfg.Version),
		Ports: []string{fmt.Sprintf("%d:3306", cfg.Port)},
		Environment: map[string]string{
			"MYSQL_DATABASE":      dbName,
			"MYSQL_USER":          "devdock",
			"MYSQL_PASSWORD":      "password",
			"MYSQL_ROOT_PASSWORD": "password",
		},
		Volumes: []string{
			fmt.Sprintf("devdock_%s_mysql_data:/var/lib/mysql", utils.NormalizeProjectName(projName)),
		},
		Networks: []string{"devdock_network"},
		Healthcheck: &Healthcheck{
			Test: []string{"CMD", "mysqladmin", "ping", "-h", "localhost"},
			Interval: "5s",
			Timeout: "5s",
			Retries: 5,
		},
	}
}

func BuildRedis(projName string, cfg config.RedisConfig) Service {
	return Service{
		Image: fmt.Sprintf("redis:%s", cfg.Version),
		Ports: []string{fmt.Sprintf("%d:6379", cfg.Port)},
		Volumes: []string{
			fmt.Sprintf("devdock_%s_redis_data:/data", utils.NormalizeProjectName(projName)),
		},
		Networks: []string{"devdock_network"},
		Healthcheck: &Healthcheck{
			Test: []string{"CMD", "redis-cli", "ping"},
			Interval: "5s",
			Timeout: "5s",
			Retries: 5,
		},
	}
}
