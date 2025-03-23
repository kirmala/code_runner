package config

type RabbitMQ struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	QueueName string `yaml:"queue_name"`
}

type HTTPConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type AppConfig struct {
	RabbitMQ `yaml:"rabbit_mq"`
	HTTPConfig `yaml:"http"`
}
