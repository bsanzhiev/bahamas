package config

type Config struct {
	Postgres PostgresConfig
	Kafka    KafkaConfig
	GRPC     GRPCConfig
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type KafkaConfig struct {
	Brokers []string
	Topic   string
	GroupID string
}

type GRPCConfig struct {
	Port int
}
