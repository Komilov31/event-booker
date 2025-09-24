package config

type Config struct {
	Postgres   PostgresConfig   `mapstructure:"postgres"`
	HttpServer HttpServerConfig `mapstructure:"http_server"`
	RabbitMq   RabbitMqConfig   `mapstructure:"rabbitmq"`
}

type PostgresConfig struct {
	Name     string `mapstructure:"name"`
	User     string `mapstructure:"user"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
}

type HttpServerConfig struct {
	Address string `mapstructure:"address"`
}

type RabbitMqConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}
