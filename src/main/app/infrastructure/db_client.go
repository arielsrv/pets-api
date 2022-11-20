package infrastructure

import (
	"context"
	"log"
	"sync"

	_ "github.com/mattn/go-sqlite3"
	"github.com/src/main/app/ent"
)

type IDbClient interface {
	Context() *ent.Client
}

type DbClient struct {
	dbMtx sync.Once
	*ent.Client
	client IDbClient
}

func NewDbClient(client IDbClient) *DbClient {
	return &DbClient{
		client: client,
	}
}

func (d *DbClient) Context() *ent.Client {
	d.dbMtx.Do(func() {
		d.Client = d.client.Context()
	})

	return d.Client
}

func (d *DbClient) Close() error {
	return d.Client.Close()
}

type MySQLClient struct {
	DbClient
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
