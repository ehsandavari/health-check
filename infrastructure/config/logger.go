package config

type SLogger struct {
	IsDevelopment     *bool  `validate:"required"`
	DisableStacktrace *bool  `validate:"required"`
	DisableStdout     *bool  `validate:"required"`
	Level             string `validate:"required"`
	Elk               SElk
	Gorm              SGorm
}

type SElk struct {
	Url           string `validate:"required"`
	TimeoutSecond byte   `validate:"required"`
}

type SGorm struct {
	SlowThresholdMilliseconds uint16 `validate:"required"`
	IgnoreRecordNotFoundError *bool  `validate:"required"`
	ParameterizedQueries      *bool  `validate:"required"`
}
