package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"health-check/infrastructure/notification"
	"health-check/infrastructure/postgres"
	"health-check/infrastructure/redis"
	"log"
	"strings"
)

type SConfig struct {
	Service      *SService             `validate:"required"`
	Jwt          *SJwt                 `validate:"required"`
	Postgres     *postgres.SConfig     `validate:"required"`
	Redis        *redis.SConfig        `validate:"required"`
	Notification *notification.SConfig `validate:"required"`
	Logger       *SLogger              `validate:"required"`
	Tracer       *STracer              `validate:"required"`
}

func NewConfig() *SConfig {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln("error in read config ", err)
	}

	config := new(SConfig)

	if err := viper.Unmarshal(config); err != nil {
		log.Fatalln("error in unmarshal config ", err)
	}

	if err := validator.New().Struct(config); err != nil {
		log.Fatalln("error in validate config ", err)
	}

	if !config.Service.Mode.IsValid() {
		log.Fatalln(config.Service.Mode.String(), "service mode is not valid !", "valid service modes is", config.Service.Mode.List())
	}

	return config
}
