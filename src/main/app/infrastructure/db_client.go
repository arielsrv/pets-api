package infrastructure

import (
	"context"
	"log"
	"sync"

	_ "github.com/mattn/go-sqlite3"
	"github.com/src/main/app/ent"
)

type IDbClient interface {
	Open() *ent.Client
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

func (d *DbClient) Open() *ent.Client {
	d.dbMtx.Do(func() {
		d.Client = d.client.Open()
	})

	return d.Client
}

func (d *DbClient) GetClient() *ent.Client {
	return d.Client
}

func (d *DbClient) Close() error {
	return d.Client.Close()
}

type MySQLClient struct {
	DbClient
	connectionString string
}

func NewMySQLClient(connectionString string) *MySQLClient {
	return &MySQLClient{connectionString: connectionString}
}

func (m *MySQLClient) Open() *ent.Client {
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
