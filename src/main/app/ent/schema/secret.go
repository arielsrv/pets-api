package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Secret holds the schema definition for the Secret entity.
type Secret struct {
	ent.Schema
}

// Fields of the Secret.
func (Secret) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").StructTag(`json:"oid,omitempty"`),
		field.String("key"),
		field.String("value"),
		field.Int64("app_id"),
		field.Bool("active").Default(true),
	}
}

// Edges of the Secret.
func (Secret) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("app", App.Type).
			Required().
			Field("app_id").
			Unique(),
	}
}
