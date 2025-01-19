package database

import (
	"context"
	"log"
	transporterqueue "netwatch/internal/pkg/transporter/queue"
	"os"

	"github.com/jackc/pgx/v5"
)

func Setup() *pgx.Conn {
	log.Println("Setting up DB")
	conn, err := pgx.Connect(context.Background(), os.Getenv("NW_DB_CONN_URL"))
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	return conn
}

func Send(item transporterqueue.QueueItem, dbConn *pgx.Conn) {
	// saves a TransportQueueItem to the db
	log.Println("Sending to DB")
}
