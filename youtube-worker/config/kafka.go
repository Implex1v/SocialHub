package config

type KafkaConfig struct {
	Host string
	Port string
}

func NewKafkaConfig() KafkaConfig {
	return KafkaConfig{
		Host: getEnvOr("KAFKA_HOST", "localhost"),
		Port: getEnvOr("KAFKA_PORT", "19092"),
	}
}
