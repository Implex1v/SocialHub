package config

type HttpConfig struct {
	Port string
}

func NewHttpConfig() *HttpConfig {
	return &HttpConfig{
		Port: getEnvOr("HTTP_PORT", "8000"),
	}
}
