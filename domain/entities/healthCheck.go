package entities

import (
	"gorm.io/datatypes"
	"health-check/domain/enums"
)

type HealthCheck struct {
	Id       uint                                  `gorm:"primaryKey;"`
	Interval string                                `gorm:"size:30;not null"`
	Url      string                                `gorm:"size:600;not null"`
	Method   enums.HttpMethod                      `gorm:"size:30;not null"`
	Headers  datatypes.JSONType[map[string]string] `gorm:"not null"`
	Body     datatypes.JSONType[map[string]any]    `gorm:"not null"`
	Status   enums.Status                          `gorm:"size:30;not null"`
	Base3
}

func NewHealthCheck(interval string, url string, method enums.HttpMethod, headers map[string]string, body map[string]any, status enums.Status) HealthCheck {
	return HealthCheck{
		Interval: interval,
		Url:      url,
		Method:   method,
		Headers:  datatypes.NewJSONType(headers),
		Body:     datatypes.NewJSONType(body),
		Status:   status,
	}
}

func (r *HealthCheck) SetStatus(status enums.Status) {
	r.Status = status
}
