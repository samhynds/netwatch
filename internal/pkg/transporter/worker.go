package transporter

import (
	"netwatch/internal/pkg/queue"
	"netwatch/internal/pkg/transporter/database"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Worker(item queue.ProcessedQueueItem, dbConnPool *pgxpool.Pool, kafkaConn string) {

	if dbConnPool != nil {
		database.Send(item, dbConnPool)
	}

	// if kafka != nil {}
}
