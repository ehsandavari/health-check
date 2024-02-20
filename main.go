package main

import (
	"github.com/ehsandavari/go-context-plus"
	"github.com/ehsandavari/go-graceful-shutdown"
	"github.com/joho/godotenv"
	"health-check/application"
	"health-check/infrastructure"
	"health-check/persistence"
	"health-check/presentation"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	loadEnv()
	run()
}

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln("error in loading .env file", err)
	}
}

func run() {
	ctx := contextplus.Background()
	newInfrastructure := infrastructure.NewInfrastructure()
	newApplication := application.NewApplication(newInfrastructure, persistence.NewPersistence(newInfrastructure))
	newApplication.StartJobs(ctx)
	newPresentation := presentation.NewPresentation(newInfrastructure, newApplication)
	newPresentation.Setup()

	newInfrastructure.ILogger.WithAny("config", newInfrastructure.SConfig).Info(ctx, "app config")

	shutdownFunc := func() {
		newPresentation.Close()
	}
	cleanupFunc := func() {
		newApplication.StopJobs(ctx)
		newInfrastructure.Close()
	}
	graceful.Shutdown(shutdownFunc, cleanupFunc, newInfrastructure.SConfig.Service.GracefulShutdownSecond)
}
