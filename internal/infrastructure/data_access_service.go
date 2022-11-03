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

type DataAccessService struct {
	dbMtx  sync.Once
	client *ent.Client
}

func NewDataAccessService() *DataAccessService {
	return &DataAccessService{}
}

func (d *DataAccessService) Open() *ent.Client {
	d.dbMtx.Do(func() {
		user := config.String("database.user")
		password := config.String("database.password")
		host := config.String("database.host")
		port := config.String("database.port")
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/pets?parseTime=True", user, password, host, port)
		dataAccess, err := ent.Open("mysql", dsn)
		if err != nil {
			log.Fatalf("failed opening connection to mysql: %v", err)
		}
		// Run the auto migration tool.
		if err = dataAccess.Schema.Create(context.Background()); err != nil {
			log.Fatalf("failed creating schema resources: %v", err)
		}
		d.client = dataAccess
	})

	return d.client
}

func (d *DataAccessService) Test(t *testing.T) *ent.Client {
	d.dbMtx.Do(func() {
		opts := []enttest.Option{
			enttest.WithOptions(ent.Debug()),
		}
		dataAccess := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1", opts...)
		d.client = dataAccess
		Seeds(d.client)
	})

	return d.client
}

func Seeds(client *ent.Client) {
	client.AppType.Create().SetID(1).SetName("backend").Save(context.Background())
	client.App.Create().SetName("customers-api").SetProjectId(2).SetAppTypeID(1).Save(context.Background())
}

func (d *DataAccessService) GetClient() *ent.Client {
	return d.client
}

func (d *DataAccessService) Close() error {
	return d.client.Close()
}
