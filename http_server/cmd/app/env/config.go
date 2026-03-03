package env

type Config struct {
	PostgresPassword string `env:"POSTGRES_PASSWORD,required"`
	PostgresUser     string `env:"POSTGRES_USER,required"`
	PostgresDB       string `env:"POSTGRES_DB,required"`

	PostgresHost string `env:"POSTGRES_HOST,required"`
	PostgresPort string `env:"POSTGRES_PORT,required"`

	RedisPassword  string   `env:"REDIS_PASSWORD,required"`
	RedisAddresses []string `env:"REDIS_ADDRESSES,required" envSeparator:","`
}
