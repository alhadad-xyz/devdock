package config

type Config struct {
	Version  string         `yaml:"version"`
	Project  ProjectConfig  `yaml:"project"`
	App      AppConfig      `yaml:"app"`
	Runtime  RuntimeConfig  `yaml:"runtime,omitempty"`
	Services ServicesConfig `yaml:"services,omitempty"`
	Commands CommandsConfig `yaml:"commands,omitempty"`
}

type ProjectConfig struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"` // laravel, nextjs, docker-compose
}

type AppConfig struct {
	Command string `yaml:"command,omitempty"`
	Port    int    `yaml:"port,omitempty"`
	RunMode string `yaml:"run_mode,omitempty"` // host or container (v0.1 only host)
}

type RuntimeConfig struct {
	// Placeholder for runtime-specific configuration if needed
}

type ServicesConfig struct {
	Postgres *PostgresConfig `yaml:"postgres,omitempty"`
	MySQL    *MySQLConfig    `yaml:"mysql,omitempty"`
	Redis    *RedisConfig    `yaml:"redis,omitempty"`
	Mailpit  *MailpitConfig  `yaml:"mailpit,omitempty"`
	MinIO    *MinIOConfig    `yaml:"minio,omitempty"`
}

type PostgresConfig struct {
	Enabled bool   `yaml:"enabled"`
	Port    int    `yaml:"port,omitempty"`
	Version string `yaml:"version,omitempty"`
}

type MySQLConfig struct {
	Enabled bool   `yaml:"enabled"`
	Port    int    `yaml:"port,omitempty"`
	Version string `yaml:"version,omitempty"`
}

type RedisConfig struct {
	Enabled bool   `yaml:"enabled"`
	Port    int    `yaml:"port,omitempty"`
	Version string `yaml:"version,omitempty"`
}

type MailpitConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Version  string `yaml:"version,omitempty"`
	SMTPPort int    `yaml:"smtp_port,omitempty"`
	UIPort   int    `yaml:"ui_port,omitempty"`
}

type MinIOConfig struct {
	Enabled     bool   `yaml:"enabled"`
	Version     string `yaml:"version,omitempty"`
	APIPort     int    `yaml:"api_port,omitempty"`
	ConsolePort int    `yaml:"console_port,omitempty"`
}

type CommandsConfig map[string]string
