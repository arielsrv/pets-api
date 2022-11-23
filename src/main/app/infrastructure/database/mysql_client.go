package database

import (
	"context"
	"log"

	"github.com/src/main/app/ent"
)

type MySQLClient struct {
	DBClient
	connectionString string
}

func NewMySQLClient(connectionString string) IDbClient {
	return &MySQLClient{
		connectionString: connectionString,
	}
}

func (m *MySQLClient) Context() *ent.Client {
	client, err := ent.Open("mysql", m.connectionString)
	if err != nil {
		log.Fatalf("failed opening connection to mysql: %v", err)
	}
	// Run the auto migration tool.
	if err = client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	return client
}
