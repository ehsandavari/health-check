package commands

type SHealthCheckDeleteCommand struct {
	id uint
}

func NewHealthCheckDeleteCommand(id uint) SHealthCheckDeleteCommand {
	return SHealthCheckDeleteCommand{
		id: id,
	}
}
