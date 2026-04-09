package config

type PostgresDB struct {
	Password string `env:"POSTGRES_PASSWORD" env-required:"true"`
	User     string `yaml:"user" env:"POSTGRES_USER"`
	DB       string `yaml:"db" env:"POSTGRES_DB"`

	Host string `yaml:"host" env:"POSTGRES_HOST"`
	Port string `yaml:"port" env:"POSTGRES_PORT"`
}

type RedisDB struct {
	Password  string   `env:"REDIS_PASSWORD" env-required:"true"`
	Addresses []string `yaml:"addresses" env:"REDIS_ADDRESSES" env-separator:","`
}

type RabbitMQ struct {
	Host      string `yaml:"host"`
	Port      string `yaml:"port"`
	QueueName string `yaml:"queue_name"`
}

type HTTPConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type PrometheusConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type AppConfig struct {
	RabbitMQ         `yaml:"rabbit_mq"`
	HTTPConfig       `yaml:"http"`
	PrometheusConfig `yaml:"prometheus"`
	PostgresDB       `yaml:"postgres_db"`
	RedisDB          `yaml:"redis_db"`
}
