package config

type RabbitMQ struct {
	Host      string `yaml:"host"`
	Port      string `yaml:"port"`
	QueueName string `yaml:"queue_name"`
}

type PostgresDB struct {
	Password string `env:"POSTGRES_PASSWORD" env-required:"true"`
	User     string `yaml:"user" env:"POSTGRES_USER"`
	DB       string `yaml:"db" env:"POSTGRES_DB"`

	Host string `yaml:"host" env:"POSTGRES_HOST"`
	Port string `yaml:"port" env:"POSTGRES_PORT"`
}

type ContainerResource struct {
	Memory    int64  `yaml:"memory"`
	NanoCPUs  int64  `yaml:"nano_cpus"`
	PidsLimit *int64 `yaml:"pids_limit"`
}

type Runner struct {
	ImageName         string            `yaml:"image_name"`
	ClientVersion     string            `yaml:"client_version"`
	ContainerResource ContainerResource `yaml:"container_resource"`
}

type AppConfig struct {
	RabbitMQ   `yaml:"rabbit_mq"`
	PostgresDB `yaml:"postgres_db"`
	Runner     `yaml:"runner"`
}
