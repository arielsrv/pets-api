// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/src/main/app/ent/app"
	"github.com/src/main/app/ent/apptype"
)

// App is the model entity for the App schema.
type App struct {
	config `json:"-"`
	// ID of the ent.
	ID int64 `json:"oid,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// ExternalGitlabProjectID holds the value of the "external_gitlab_project_id" field.
	ExternalGitlabProjectID int64 `json:"external_gitlab_project_id,omitempty"`
	// AppTypeID holds the value of the "app_type_id" field.
	AppTypeID int `json:"app_type_id,omitempty"`
	// Active holds the value of the "active" field.
	Active bool `json:"active,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the AppQuery when eager-loading is set.
	Edges AppEdges `json:"edges"`
}

// AppEdges holds the relations/edges for other nodes in the graph.
type AppEdges struct {
	// AppsTypes holds the value of the apps_types edge.
	AppsTypes *AppType `json:"apps_types,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// AppsTypesOrErr returns the AppsTypes value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e AppEdges) AppsTypesOrErr() (*AppType, error) {
	if e.loadedTypes[0] {
		if e.AppsTypes == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: apptype.Label}
		}
		return e.AppsTypes, nil
	}
	return nil, &NotLoadedError{edge: "apps_types"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*App) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case app.FieldActive:
			values[i] = new(sql.NullBool)
		case app.FieldID, app.FieldExternalGitlabProjectID, app.FieldAppTypeID:
			values[i] = new(sql.NullInt64)
		case app.FieldName:
			values[i] = new(sql.NullString)
		default:
			return nil, fmt.Errorf("unexpected column %q for type App", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the App fields.
func (a *App) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case app.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			a.ID = int64(value.Int64)
		case app.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				a.Name = value.String
			}
		case app.FieldExternalGitlabProjectID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field external_gitlab_project_id", values[i])
			} else if value.Valid {
				a.ExternalGitlabProjectID = value.Int64
			}
		case app.FieldAppTypeID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field app_type_id", values[i])
			} else if value.Valid {
				a.AppTypeID = int(value.Int64)
			}
		case app.FieldActive:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field active", values[i])
			} else if value.Valid {
				a.Active = value.Bool
			}
		}
	}
	return nil
}

// QueryAppsTypes queries the "apps_types" edge of the App entity.
func (a *App) QueryAppsTypes() *AppTypeQuery {
	return (&AppClient{config: a.config}).QueryAppsTypes(a)
}

// Update returns a builder for updating this App.
// Note that you need to call App.Unwrap() before calling this method if this App
// was returned from a transaction, and the transaction was committed or rolled back.
func (a *App) Update() *AppUpdateOne {
	return (&AppClient{config: a.config}).UpdateOne(a)
}

// Unwrap unwraps the App entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (a *App) Unwrap() *App {
	_tx, ok := a.config.driver.(*txDriver)
	if !ok {
		panic("ent: App is not a transactional entity")
	}
	a.config.driver = _tx.drv
	return a
}

// String implements the fmt.Stringer.
func (a *App) String() string {
	var builder strings.Builder
	builder.WriteString("App(")
	builder.WriteString(fmt.Sprintf("id=%v, ", a.ID))
	builder.WriteString("name=")
	builder.WriteString(a.Name)
	builder.WriteString(", ")
	builder.WriteString("external_gitlab_project_id=")
	builder.WriteString(fmt.Sprintf("%v", a.ExternalGitlabProjectID))
	builder.WriteString(", ")
	builder.WriteString("app_type_id=")
	builder.WriteString(fmt.Sprintf("%v", a.AppTypeID))
	builder.WriteString(", ")
	builder.WriteString("active=")
	builder.WriteString(fmt.Sprintf("%v", a.Active))
	builder.WriteByte(')')
	return builder.String()
}

// Apps is a parsable slice of App.
type Apps []*App

func (a Apps) config(cfg config) {
	for _i := range a {
		a[_i].config = cfg
	}
}
