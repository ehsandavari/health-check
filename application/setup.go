package application

import (
	"github.com/ehsandavari/go-context-plus"
	"health-check/application/handlers/commands"
	"health-check/application/handlers/jobs"
	"health-check/application/handlers/queries"
	"health-check/infrastructure"
	"health-check/persistence"
)

type Application struct {
	infrastructure *infrastructure.Infrastructure
	Commands       commands.Commands
	Queries        queries.Queries
	Jobs           jobs.Jobs
}

func NewApplication(infrastructure *infrastructure.Infrastructure, persistence *persistence.Persistence) *Application {
	return &Application{
		infrastructure: infrastructure,
		Commands:       commands.NewCommands(infrastructure, persistence),
		Queries:        queries.NewQueries(infrastructure, persistence),
		Jobs:           jobs.NewJobs(infrastructure, persistence),
	}
}

func (r Application) StartJobs(ctx *contextplus.Context) {
	if err := r.Jobs.HealthCheck.Start(ctx); err != nil {
		r.infrastructure.ILogger.WithError(err).Fatal(ctx, "error in start health check job")
	}
}

func (r Application) StopJobs(ctx *contextplus.Context) {
	if err := r.Jobs.HealthCheck.Stop(ctx); err != nil {
		r.infrastructure.ILogger.WithError(err).Error(ctx, "error in stop health check job")
	}
}
