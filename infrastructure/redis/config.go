package redis

type SConfig struct {
	Host string `validate:"required"`
	Port string `validate:"required"`
}
