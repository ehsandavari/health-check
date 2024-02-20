package commands

import "health-check/domain/enums"

type SHealthCheckCreateCommand struct {
	interval string
	url      string
	method   enums.HttpMethod
	headers  map[string]string
	body     map[string]any
}

func NewHealthCheckCreateCommand(interval string, url string, method enums.HttpMethod, headers map[string]string, body map[string]any) SHealthCheckCreateCommand {
	return SHealthCheckCreateCommand{
		interval: interval,
		url:      url,
		method:   method,
		headers:  headers,
		body:     body,
	}
}
