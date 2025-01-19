package transporter

import (
	"netwatch/internal/pkg/transporter/database"
	transporterqueue "netwatch/internal/pkg/transporter/queue"

	"github.com/jackc/pgx/v5"
)

func Worker(item transporterqueue.QueueItem, dbConn *pgx.Conn, kafkaConn string) {

	if dbConn != nil {
		database.Send(item, dbConn)
	}

	// if kafka != nil {}
}
