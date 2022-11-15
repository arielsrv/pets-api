package infrastructure

import (
	"context"
	"fmt"
	"log"
	"sync"
	"testing"

	"github.com/internal/config"

	"github.com/ent"
	"github.com/ent/enttest"
	_ "github.com/mattn/go-sqlite3"
)

type DbClient struct {
	dbMtx sync.Once
	*ent.Client
}

func NewDbClient() *DbClient {
	return &DbClient{}
}

func (d *DbClient) Open() *ent.Client {
	d.dbMtx.Do(func() {
		user := config.String("database.user")
		password := config.String("database.password")
		host := config.String("database.host")
		port := config.String("database.port")
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/pets?parseTime=True", user, password, host, port)
		dbClient, err := ent.Open("mysql", dsn)
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
