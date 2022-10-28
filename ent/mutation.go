// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/ent/app"
	"github.com/ent/predicate"

	"entgo.io/ent"
)

const (
	// Operation types.
	OpCreate    = ent.OpCreate
	OpDelete    = ent.OpDelete
	OpDeleteOne = ent.OpDeleteOne
	OpUpdate    = ent.OpUpdate
	OpUpdateOne = ent.OpUpdateOne

	// Node types.
	TypeApp = "App"
)

// AppMutation represents an operation that mutates the App nodes in the graph.
type AppMutation struct {
	config
	op            Op
	typ           string
	id            *int64
	name          *string
	projectId     *int64
	addprojectId  *int64
	active        *bool
	clearedFields map[string]struct{}
	done          bool
	oldValue      func(context.Context) (*App, error)
	predicates    []predicate.App
}

var _ ent.Mutation = (*AppMutation)(nil)

// appOption allows management of the mutation configuration using functional options.
type appOption func(*AppMutation)

// newAppMutation creates new mutation for the App entity.
func newAppMutation(c config, op Op, opts ...appOption) *AppMutation {
	m := &AppMutation{
		config:        c,
		op:            op,
		typ:           TypeApp,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withAppID sets the ID field of the mutation.
func withAppID(id int64) appOption {
	return func(m *AppMutation) {
		var (
			err   error
			once  sync.Once
			value *App
		)
		m.oldValue = func(ctx context.Context) (*App, error) {
			once.Do(func() {
				if m.done {
					err = errors.New("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().App.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withApp sets the old App of the mutation.
func withApp(node *App) appOption {
	return func(m *AppMutation) {
		m.oldValue = func(context.Context) (*App, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m AppMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m AppMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, errors.New("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// SetID sets the value of the id field. Note that this
// operation is only accepted on creation of App entities.
func (m *AppMutation) SetID(id int64) {
	m.id = &id
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *AppMutation) ID() (id int64, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *AppMutation) IDs(ctx context.Context) ([]int64, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []int64{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().App.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetName sets the "name" field.
func (m *AppMutation) SetName(s string) {
	m.name = &s
}

// Name returns the value of the "name" field in the mutation.
func (m *AppMutation) Name() (r string, exists bool) {
	v := m.name
	if v == nil {
		return
	}
	return *v, true
}

// OldName returns the old "name" field's value of the App entity.
// If the App object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *AppMutation) OldName(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldName is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldName requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldName: %w", err)
	}
	return oldValue.Name, nil
}

// ResetName resets all changes to the "name" field.
func (m *AppMutation) ResetName() {
	m.name = nil
}

// SetProjectId sets the "projectId" field.
func (m *AppMutation) SetProjectId(i int64) {
	m.projectId = &i
	m.addprojectId = nil
}

// ProjectId returns the value of the "projectId" field in the mutation.
func (m *AppMutation) ProjectId() (r int64, exists bool) {
	v := m.projectId
	if v == nil {
		return
	}
	return *v, true
}

// OldProjectId returns the old "projectId" field's value of the App entity.
// If the App object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *AppMutation) OldProjectId(ctx context.Context) (v int64, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldProjectId is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldProjectId requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldProjectId: %w", err)
	}
	return oldValue.ProjectId, nil
}

// AddProjectId adds i to the "projectId" field.
func (m *AppMutation) AddProjectId(i int64) {
	if m.addprojectId != nil {
		*m.addprojectId += i
	} else {
		m.addprojectId = &i
	}
}

// AddedProjectId returns the value that was added to the "projectId" field in this mutation.
func (m *AppMutation) AddedProjectId() (r int64, exists bool) {
	v := m.addprojectId
	if v == nil {
		return
	}
	return *v, true
}

// ResetProjectId resets all changes to the "projectId" field.
func (m *AppMutation) ResetProjectId() {
	m.projectId = nil
	m.addprojectId = nil
}

// SetActive sets the "active" field.
func (m *AppMutation) SetActive(b bool) {
	m.active = &b
}

// Active returns the value of the "active" field in the mutation.
func (m *AppMutation) Active() (r bool, exists bool) {
	v := m.active
	if v == nil {
		return
	}
	return *v, true
}

// OldActive returns the old "active" field's value of the App entity.
// If the App object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *AppMutation) OldActive(ctx context.Context) (v bool, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldActive is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldActive requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldActive: %w", err)
	}
	return oldValue.Active, nil
}

// ResetActive resets all changes to the "active" field.
func (m *AppMutation) ResetActive() {
	m.active = nil
}

// Where appends a list predicates to the AppMutation builder.
func (m *AppMutation) Where(ps ...predicate.App) {
	m.predicates = append(m.predicates, ps...)
}

// Op returns the operation name.
func (m *AppMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (App).
func (m *AppMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *AppMutation) Fields() []string {
	fields := make([]string, 0, 3)
	if m.name != nil {
		fields = append(fields, app.FieldName)
	}
	if m.projectId != nil {
		fields = append(fields, app.FieldProjectId)
	}
	if m.active != nil {
		fields = append(fields, app.FieldActive)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *AppMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case app.FieldName:
		return m.Name()
	case app.FieldProjectId:
		return m.ProjectId()
	case app.FieldActive:
		return m.Active()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *AppMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case app.FieldName:
		return m.OldName(ctx)
	case app.FieldProjectId:
		return m.OldProjectId(ctx)
	case app.FieldActive:
		return m.OldActive(ctx)
	}
	return nil, fmt.Errorf("unknown App field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *AppMutation) SetField(name string, value ent.Value) error {
	switch name {
	case app.FieldName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetName(v)
		return nil
	case app.FieldProjectId:
		v, ok := value.(int64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetProjectId(v)
		return nil
	case app.FieldActive:
		v, ok := value.(bool)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetActive(v)
		return nil
	}
	return fmt.Errorf("unknown App field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *AppMutation) AddedFields() []string {
	var fields []string
	if m.addprojectId != nil {
		fields = append(fields, app.FieldProjectId)
	}
	return fields
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *AppMutation) AddedField(name string) (ent.Value, bool) {
	switch name {
	case app.FieldProjectId:
		return m.AddedProjectId()
	}
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *AppMutation) AddField(name string, value ent.Value) error {
	switch name {
	case app.FieldProjectId:
		v, ok := value.(int64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddProjectId(v)
		return nil
	}
	return fmt.Errorf("unknown App numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *AppMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *AppMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *AppMutation) ClearField(name string) error {
	return fmt.Errorf("unknown App nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *AppMutation) ResetField(name string) error {
	switch name {
	case app.FieldName:
		m.ResetName()
		return nil
	case app.FieldProjectId:
		m.ResetProjectId()
		return nil
	case app.FieldActive:
		m.ResetActive()
		return nil
	}
	return fmt.Errorf("unknown App field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *AppMutation) AddedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *AppMutation) AddedIDs(name string) []ent.Value {
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *AppMutation) RemovedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *AppMutation) RemovedIDs(name string) []ent.Value {
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *AppMutation) ClearedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *AppMutation) EdgeCleared(name string) bool {
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *AppMutation) ClearEdge(name string) error {
	return fmt.Errorf("unknown App unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *AppMutation) ResetEdge(name string) error {
	return fmt.Errorf("unknown App edge %s", name)
}
