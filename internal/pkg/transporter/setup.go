package transporter

import (
	"netwatch/internal/pkg/config"
	"netwatch/internal/pkg/transporter/database"

	"github.com/jackc/pgx/v5"
)

func SetupConnections(config *config.Config) (*pgx.Conn, string) {
	var dbConnection *pgx.Conn
	var kafkaConnection string

	if config.Config.Transporters.Database {
		dbConnection = database.Setup()
	}

	if config.Config.Transporters.Queue {
		// kafkaConnection, err := kafka.Setup()
	}

	return dbConnection, kafkaConnection
}
