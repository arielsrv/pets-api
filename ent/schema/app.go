package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// App holds the schema definition for the App entity.
type App struct {
	ent.Schema
}

// Fields of the App.
func (App) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			StructTag(`json:"oid,omitempty"`),
		field.String("name").
			Unique(),
		field.Int64("projectId").
			Unique(),
		field.Bool("active").
			Default(true),
	}
}

// Edges of the App.
func (App) Edges() []ent.Edge {
	return nil
}
