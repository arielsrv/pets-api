package infrastructure

import (
	"testing"

	"github.com/src/main/app/ent"
	"github.com/src/main/app/ent/enttest"
)

type SQLiteClient struct {
	DbClient
	t *testing.T
}

func NewSQLiteClient(t *testing.T) IDbClient {
	return &SQLiteClient{t: t}
}

func (s *SQLiteClient) Context() *ent.Client {
	opts := []enttest.Option{
		enttest.WithOptions(ent.Debug()),
	}

	return enttest.Open(s.t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1", opts...)
}
