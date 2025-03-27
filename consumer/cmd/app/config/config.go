package config

type RabbitMQ struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	QueueName string `yaml:"queue_name"`
}

type Repository struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	Address string `yaml:"address"`
} 

type CodeProcessor struct {
	ImageName string `yaml:"image_name"`
}

type AppConfig struct {
	RabbitMQ `yaml:"rabbit_mq"`
	Repository `yaml:"repository"`
	CodeProcessor `yaml:"code_processor"`
}
