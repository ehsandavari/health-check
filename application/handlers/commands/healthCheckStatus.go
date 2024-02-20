package commands

import "health-check/domain/enums"

type SHealthCheckStatusCommand struct {
	id     uint
	status enums.Status
}

func NewHealthCheckStatusCommand(id uint, status enums.Status) SHealthCheckStatusCommand {
	return SHealthCheckStatusCommand{
		id:     id,
		status: status,
	}
}
