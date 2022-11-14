// Code generated by ent, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// AppsColumns holds the columns for the "apps" table.
	AppsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt64, Increment: true},
		{Name: "name", Type: field.TypeString, Unique: true},
		{Name: "project_id", Type: field.TypeInt64, Unique: true},
		{Name: "active", Type: field.TypeBool, Default: true},
		{Name: "app_type_id", Type: field.TypeInt},
	}
	// AppsTable holds the schema information for the "apps" table.
	AppsTable = &schema.Table{
		Name:       "apps",
		Columns:    AppsColumns,
		PrimaryKey: []*schema.Column{AppsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "apps_apps_types_apps",
				Columns:    []*schema.Column{AppsColumns[4]},
				RefColumns: []*schema.Column{AppsTypesColumns[0]},
				OnDelete:   schema.NoAction,
			},
		},
	}
	// AppsTypesColumns holds the columns for the "apps_types" table.
	AppsTypesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "name", Type: field.TypeString, Unique: true},
	}
	// AppsTypesTable holds the schema information for the "apps_types" table.
	AppsTypesTable = &schema.Table{
		Name:       "apps_types",
		Columns:    AppsTypesColumns,
		PrimaryKey: []*schema.Column{AppsTypesColumns[0]},
	}
	// SecretsColumns holds the columns for the "secrets" table.
	SecretsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt64, Increment: true},
		{Name: "key", Type: field.TypeString},
		{Name: "value", Type: field.TypeString},
		{Name: "active", Type: field.TypeBool, Default: true},
		{Name: "app_id", Type: field.TypeInt64},
	}
	// SecretsTable holds the schema information for the "secrets" table.
	SecretsTable = &schema.Table{
		Name:       "secrets",
		Columns:    SecretsColumns,
		PrimaryKey: []*schema.Column{SecretsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "secrets_apps_app",
				Columns:    []*schema.Column{SecretsColumns[4]},
				RefColumns: []*schema.Column{AppsColumns[0]},
				OnDelete:   schema.NoAction,
			},
		},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		AppsTable,
		AppsTypesTable,
		SecretsTable,
	}
)

func init() {
	AppsTable.ForeignKeys[0].RefTable = AppsTypesTable
	AppsTypesTable.Annotation = &entsql.Annotation{
		Table: "apps_types",
	}
	SecretsTable.ForeignKeys[0].RefTable = AppsTable
}
