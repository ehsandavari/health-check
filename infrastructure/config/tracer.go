package config

type STracer struct {
	Host string `validate:"required"`
	Port string `validate:"required"`
}
