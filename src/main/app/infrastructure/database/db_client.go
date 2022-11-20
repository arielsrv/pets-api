package database

import (
	"sync"

	_ "github.com/mattn/go-sqlite3"
	"github.com/src/main/app/ent"
)

type IDbClient interface {
	Context() *ent.Client
}

type DbClient struct {
	dbClient IDbClient
	mutex    sync.Once
	*ent.Client
}

func NewDbClient(client IDbClient) *DbClient {
	return &DbClient{
		dbClient: client,
	}
}

// Context template method, used by concrete impl.
func (d *DbClient) Context() *ent.Client {
	d.mutex.Do(func() {
		d.Client = d.dbClient.Context()
	})

	return d.Client
}

func (d *DbClient) Close() error {
	return d.Client.Close()
}
