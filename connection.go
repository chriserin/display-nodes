package main

import (
	"context"
	"fmt"
	"os"

	pgx "github.com/jackc/pgx/v5"
)

type Connection struct {
	databaseUrl string
	conn        *pgx.Conn
}

func (c *Connection) Connect() {
	conn, err := pgx.Connect(context.Background(), c.databaseUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	c.conn = conn
}

func (c Connection) ExecuteExplain(query string) string {

	var explainResult string
	err := c.conn.QueryRow(context.Background(), query).Scan(&explainResult)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	return explainResult
}

func (c Connection) Close() {
	c.conn.Close(context.Background())
}
