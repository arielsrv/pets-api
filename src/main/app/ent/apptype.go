// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/arielsrv/pets-api/src/main/app/ent/apptype"
)

// AppType is the model entity for the AppType schema.
type AppType struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"oid,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the AppTypeQuery when eager-loading is set.
	Edges        AppTypeEdges `json:"edges"`
	selectValues sql.SelectValues
}

// AppTypeEdges holds the relations/edges for other nodes in the graph.
type AppTypeEdges struct {
	// Apps holds the value of the apps edge.
	Apps []*App `json:"apps,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// AppsOrErr returns the Apps value or an error if the edge
// was not loaded in eager-loading.
func (e AppTypeEdges) AppsOrErr() ([]*App, error) {
	if e.loadedTypes[0] {
		return e.Apps, nil
	}
	return nil, &NotLoadedError{edge: "apps"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*AppType) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case apptype.FieldID:
			values[i] = new(sql.NullInt64)
		case apptype.FieldName:
			values[i] = new(sql.NullString)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the AppType fields.
func (at *AppType) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case apptype.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			at.ID = int(value.Int64)
		case apptype.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				at.Name = value.String
			}
		default:
			at.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the AppType.
// This includes values selected through modifiers, order, etc.
func (at *AppType) Value(name string) (ent.Value, error) {
	return at.selectValues.Get(name)
}

// QueryApps queries the "apps" edge of the AppType entity.
func (at *AppType) QueryApps() *AppQuery {
	return NewAppTypeClient(at.config).QueryApps(at)
}

// Update returns a builder for updating this AppType.
// Note that you need to call AppType.Unwrap() before calling this method if this AppType
// was returned from a transaction, and the transaction was committed or rolled back.
func (at *AppType) Update() *AppTypeUpdateOne {
	return NewAppTypeClient(at.config).UpdateOne(at)
}

// Unwrap unwraps the AppType entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (at *AppType) Unwrap() *AppType {
	_tx, ok := at.config.driver.(*txDriver)
	if !ok {
		panic("ent: AppType is not a transactional entity")
	}
	at.config.driver = _tx.drv
	return at
}

// String implements the fmt.Stringer.
func (at *AppType) String() string {
	var builder strings.Builder
	builder.WriteString("AppType(")
	builder.WriteString(fmt.Sprintf("id=%v, ", at.ID))
	builder.WriteString("name=")
	builder.WriteString(at.Name)
	builder.WriteByte(')')
	return builder.String()
}

// AppTypes is a parsable slice of AppType.
type AppTypes []*AppType
