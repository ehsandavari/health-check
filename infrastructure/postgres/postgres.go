package postgres

import (
	"github.com/ehsandavari/go-context-plus"
	"github.com/ehsandavari/go-logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"health-check/domain/entities"
)

type SPostgres struct {
	Database *gorm.DB
}

func NewPostgres(config *SConfig, logger logger.ILogger) SPostgres {
	sPostgres := new(SPostgres)
	var err error

	if sPostgres.Database, err = gorm.Open(
		postgres.Open("host="+config.Host+" user="+config.User+" password="+config.Password+" dbname="+config.DatabaseName+" port="+config.Port+" sslmode="+config.SslMode+" TimeZone="+config.TimeZone+""),
		&gorm.Config{
			Logger: logger.GormLogger(),
		},
	); err != nil {
		logger.WithError(err).Fatal(contextplus.Background(), "error in connect to postgres")
	}

	if err = sPostgres.setup(); err != nil {
		logger.WithError(err).Fatal(contextplus.Background(), "error in setup postgres")
	}

	return *sPostgres
}

func (r *SPostgres) setup() error {
	return r.Database.AutoMigrate(
		new(entities.HealthCheck),
		new(entities.HealthCheckRequest),
	)
}

func (r *SPostgres) Close() error {
	db, err := r.Database.DB()
	if err != nil {
		return err
	}
	return db.Close()
}
