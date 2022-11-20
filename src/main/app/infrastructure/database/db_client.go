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
