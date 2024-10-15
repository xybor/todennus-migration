package postgres

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	config "github.com/xybor/todennus-config"
	"github.com/xybor/x/xcontext"
	postgresDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func Initialize(ctx context.Context, config config.Config) (*gorm.DB, error) {
	loglevel := config.Variable.Postgres.LogLevel
	newLogger := gormlogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		gormlogger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  gormlogger.LogLevel(loglevel),
			IgnoreRecordNotFoundError: true,
		},
	)

	dsnFormat := "host=%s user=%s password=%s dbname=%s port=%d sslmode=%s"
	host := config.Variable.Postgres.Host
	port := config.Variable.Postgres.Port
	user := config.Secret.Postgres.User
	password := config.Secret.Postgres.Password
	dbname := config.Secret.Postgres.DBName
	sslmode := config.Variable.Postgres.SSLMode
	dsn := fmt.Sprintf(dsnFormat, host, user, password, dbname, port, sslmode)

	timezone := config.Variable.Postgres.Timezone
	if timezone != "" {
		dsn += fmt.Sprintf(" TimeZone=%s", timezone)
	}

	var postgresDB *gorm.DB
	var err error
	nRetries := config.Variable.Postgres.RetryAttempts
	retryDuration := config.Variable.Postgres.RetryInterval
	for i := 0; i < nRetries; i++ {
		postgresDB, err = gorm.Open(postgresDriver.Open(dsn), &gorm.Config{Logger: newLogger})
		if err == nil {
			break
		}
		xcontext.Logger(ctx).Warn("failed-to-connect-to-postgres", "err", err)
		time.Sleep(time.Duration(retryDuration) * time.Second)
	}

	if err != nil {
		return nil, err
	}

	xcontext.Logger(ctx).Info("connect postgres successfully")
	return postgresDB, nil
}
