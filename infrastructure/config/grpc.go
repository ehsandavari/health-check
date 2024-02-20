package config

type Grpc struct {
	IsEnabled     *bool  `validate:"required"`
	Port          string `validate:"required"`
	IsDevelopment *bool  `validate:"required"`
}
