package transporter

import (
	"netwatch/internal/pkg/transporter/database"
	transporterqueue "netwatch/internal/pkg/transporter/queue"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Worker(item transporterqueue.QueueItem, dbConnPool *pgxpool.Pool, kafkaConn string) {

	if dbConnPool != nil {
		database.Send(item, dbConnPool)
	}

	// if kafka != nil {}
}
