package database

import (
	"context"
	"fmt"
	"log"
	"netwatch/internal/pkg/queue"
	"os"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Setup() *pgxpool.Pool {
	dbConnPool, err := pgxpool.New(context.Background(), os.Getenv("NW_DB_CONN_URL"))
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	// Think this will close before its used by the rest of the app
	// defer conn.Close(context.Background())

	err = setupTables(context.Background(), dbConnPool)
	if err != nil {
		log.Printf("Unable to setup tables: %v\n", err)
		os.Exit(1)
	}

	return dbConnPool
}

func Send(item queue.ProcessedQueueItem, dbConnPool *pgxpool.Pool) {
	// saves a TransportQueueItem to the db
	log.Println("Sending to DB")
	id := uuid.New()
	query := `
		INSERT INTO crawl_index (id, url, links, content, html, headers, timestamp)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	if _, err := dbConnPool.Exec(context.Background(), query,
		id,
		item.URL,
		item.Links,
		item.Content,
		item.Body,
		item.Headers,
		item.Timestamp,
	); err != nil {
		log.Printf("Failed to insert into crawl_index: %v", err)
	}
}

func setupTables(ctx context.Context, dbConnPool *pgxpool.Pool) error {
	sqlFile, err := os.ReadFile("scripts/init_db.sql")
	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}

	// Execute the SQL statements
	_, err = dbConnPool.Exec(ctx, string(sqlFile))
	if err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	return nil
}
