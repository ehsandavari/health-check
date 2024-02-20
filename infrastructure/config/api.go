package config

type SApi struct {
	IsEnabled *bool  `validate:"required"`
	Mode      string `validate:"required"`
	Host      string `validate:"required"`
	Port      string `validate:"required"`
}
