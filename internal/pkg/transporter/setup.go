package transporter

import (
	"netwatch/internal/pkg/config"
	"netwatch/internal/pkg/transporter/database"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupConnections(config *config.Config) (*pgxpool.Pool, string) {
	var dbConnPool *pgxpool.Pool
	var kafkaConnection string

	if config.Config.Transporters.Database {
		dbConnPool = database.Setup()
	}

	if config.Config.Transporters.Queue {
		// kafkaConnection, err := kafka.Setup()

	}

	return dbConnPool, kafkaConnection
}
