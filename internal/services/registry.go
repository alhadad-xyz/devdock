package services

type PortDefinition struct {
	Name          string
	ContainerPort int
	DefaultPort   int
	IsPrimary     bool
}

type WebUIDefinition struct {
	PortName string
	Label    string
}

type HealthCheckDefinition struct {
	Type        string // e.g. "cmd", "http", "exec"
	Command     []string
	Endpoint    string
	Interval    string
	Timeout     string
	Retries     int
	StartPeriod string
}

type ServiceDefinition struct {
	Name          string
	Image         string
	Ports         []PortDefinition
	HasVolume     bool
	WebUI         *WebUIDefinition
	HealthCheck   HealthCheckDefinition
	InitContainer *ServiceDefinition
}

var registry = map[string]ServiceDefinition{
	"postgres": {
		Name:  "postgres",
		Image: "postgres",
		Ports: []PortDefinition{
			{Name: "db", ContainerPort: 5432, DefaultPort: 5432, IsPrimary: true},
		},
		HasVolume: true,
	},
	"mysql": {
		Name:  "mysql",
		Image: "mysql",
		Ports: []PortDefinition{
			{Name: "db", ContainerPort: 3306, DefaultPort: 3306, IsPrimary: true},
		},
		HasVolume: true,
	},
	"redis": {
		Name:  "redis",
		Image: "redis",
		Ports: []PortDefinition{
			{Name: "db", ContainerPort: 6379, DefaultPort: 6379, IsPrimary: true},
		},
		HasVolume: true,
	},
	"mailpit": {
		Name: "mailpit",
		Ports: []PortDefinition{
			{Name: "smtp", DefaultPort: 1025, IsPrimary: true},
			{Name: "ui", DefaultPort: 8025},
		},
		HasVolume: false,
		WebUI: &WebUIDefinition{
			PortName: "ui",
			Label:    "Web UI",
		},
	},
	"minio": {
		Name: "minio",
		Ports: []PortDefinition{
			{Name: "api", DefaultPort: 9000, IsPrimary: true},
			{Name: "console", DefaultPort: 9001},
		},
		HasVolume: true,
		WebUI: &WebUIDefinition{
			PortName: "console",
			Label:    "MinIO Console",
		},
	},
}

type Registry struct{}

func Get(name string) (ServiceDefinition, bool) {
	s, ok := registry[name]
	return s, ok
}

func All() []ServiceDefinition {
	var all []ServiceDefinition
	// Return in a stable order if needed, but for now slice is fine.
	for _, s := range registry {
		all = append(all, s)
	}
	return all
}
