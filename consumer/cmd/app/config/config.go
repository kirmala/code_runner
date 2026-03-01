package config

type RabbitMQ struct {
	Host      string `yaml:"host"`
	Port      string `yaml:"port"`
	QueueName string `yaml:"queue_name"`
}

type Repository struct {
	Host    string `yaml:"host"`
	Port    string `yaml:"port"`
	Address string `yaml:"address"`
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

type Postgres struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type AppConfig struct {
	RabbitMQ   `yaml:"rabbit_mq"`
	Repository `yaml:"repository"`
	Runner     `yaml:"runner"`
	Postgres   `yaml:"postgres"`
}
