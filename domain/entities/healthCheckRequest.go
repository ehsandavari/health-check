package entities

import (
	"gorm.io/datatypes"
)

type HealthCheckRequest struct {
	Id            uint                                    `gorm:"primaryKey;"`
	HealthCheckId uint                                    `gorm:"not null"`
	Headers       datatypes.JSONType[map[string][]string] `gorm:"not null"`
	Body          string                                  `gorm:"not null"`
	StatusCode    int                                     `gorm:"not null"`
	Base1

	HealthCheck HealthCheck
}

func NewHealthCheckRequest(healthCheckId uint, headers map[string][]string, body string, statusCode int) HealthCheckRequest {
	return HealthCheckRequest{
		HealthCheckId: healthCheckId,
		Headers:       datatypes.NewJSONType(headers),
		Body:          body,
		StatusCode:    statusCode,
	}
}
