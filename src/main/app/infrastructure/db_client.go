package infrastructure

import (
	"context"
	"log"
	"sync"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/src/main/app/ent"
	"github.com/src/main/app/ent/enttest"
)

type DbClient struct {
	dbMtx sync.Once
	*ent.Client
	connectionString string
}

func NewDbClient(connectionString string) *DbClient {
	return &DbClient{
		connectionString: connectionString,
	}
}

func (d *DbClient) Open() *ent.Client {
	d.dbMtx.Do(func() {
		dbClient, err := ent.Open("mysql", d.connectionString)
		if err != nil {
			log.Fatalf("failed opening connection to mysql: %v", err)
		}
		// Run the auto migration tool.
		if err = dbClient.Schema.Create(context.Background()); err != nil {
			log.Fatalf("failed creating schema resources: %v", err)
		}
		d.Client = dbClient
	})

	return d.Client
}

func (d *DbClient) Test(t *testing.T) *ent.Client {
	d.dbMtx.Do(func() {
		opts := []enttest.Option{
			enttest.WithOptions(ent.Debug()),
		}
		dbClient := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1", opts...)
		d.Client = dbClient
		Seeds(d.Client)
	})

	return d.Client
}

func Seeds(client *ent.Client) {
	client.AppType.Create().SetID(1).SetName("backend").Save(context.Background())
	client.App.Create().SetName("customers-api").SetProjectId(2).SetAppTypeID(1).Save(context.Background())
	client.Secret.Create().SetKey("PETS_CUSTOMERS-API_MYSECRETKEY").SetValue("MYSECRETVALUE").SetAppID(1).Save(context.Background())
}

func (d *DbClient) GetClient() *ent.Client {
	return d.Client
}

func (d *DbClient) Close() error {
	return d.Client.Close()
}
