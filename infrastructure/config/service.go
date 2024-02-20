package config

import (
	"health-check/domain/enums"
)

type SService struct {
	Id                     int               `validate:"required"`
	Name                   string            `validate:"required"`
	Namespace              string            `validate:"required"`
	InstanceId             string            `validate:"required"`
	Version                string            `validate:"required"`
	Mode                   enums.ServiceMode `validate:"required"`
	CommitId               string            `validate:"required"`
	GracefulShutdownSecond byte              `validate:"required"`
	Api                    *SApi             `validate:"required"`
	Grpc                   *Grpc             `validate:"required"`
}
