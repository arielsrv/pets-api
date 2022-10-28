package infrastructure

import (
	"context"
	"fmt"
	"github.com/ent"
	"github.com/ent/enttest"
	"github.com/internal/server"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"sync"
	"testing"
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
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/pets?parseTime=True",
			server.GetAppConfig().Database.User,
			server.GetAppConfig().Database.Password,
			server.GetAppConfig().Database.Host,
			server.GetAppConfig().Database.Port)
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
	})

	return d.client
}

func (d *DataAccessService) GetClient() *ent.Client {
	return d.client
}

func (d *DataAccessService) Close() error {
	return d.client.Close()
}
