package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// AppType holds the schema definition for the AppType entity.
type AppType struct {
	ent.Schema
}

// Annotations of the User.
func (AppType) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "apps_types"},
	}
}

// Fields of the AppType.
func (AppType) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			StructTag(`json:"oid,omitempty"`),
		field.String("name").
			Unique(),
	}
}

// Edges of the AppType.
func (AppType) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("apps", App.Type),
	}
}
