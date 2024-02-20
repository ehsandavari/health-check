package dtos

import (
	"health-check/domain/enums"
	"time"
)

type HealthCheckCreateRequest struct {
	Interval string            `binding:"required" example:"1h30m10s"`
	Url      string            `binding:"required,http_url" example:"https://google.com/"`
	Method   enums.HttpMethod  `binding:"required,enum"`
	Headers  map[string]string `binding:"required"`
	Body     map[string]any    `binding:"required"`
}

type HealthCheckCreateResponse struct {
	Id        uint
	Interval  string
	Url       string
	Method    enums.HttpMethod
	Headers   map[string]string
	Body      map[string]any
	Status    enums.Status
	CreatedAt time.Time
}

type HealthCheckStatusRequest struct {
	Id     uint         `binding:"required"`
	Status enums.Status `binding:"required,enum"`
}

type HealthCheckStatusResponse struct {
	Id        uint
	Status    enums.Status
	UpdatedAt time.Time
}

type HealthCheckDeleteRequest struct {
	Id uint `binding:"required"`
}

type HealthCheckDeleteResponse struct {
	Id uint
}
