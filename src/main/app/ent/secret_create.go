// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/arielsrv/pets-api/src/main/app/ent/app"
	"github.com/arielsrv/pets-api/src/main/app/ent/secret"
)

// SecretCreate is the builder for creating a Secret entity.
type SecretCreate struct {
	config
	mutation *SecretMutation
	hooks    []Hook
}

// SetKey sets the "key" field.
func (sc *SecretCreate) SetKey(s string) *SecretCreate {
	sc.mutation.SetKey(s)
	return sc
}

// SetValue sets the "value" field.
func (sc *SecretCreate) SetValue(s string) *SecretCreate {
	sc.mutation.SetValue(s)
	return sc
}

// SetAppID sets the "app_id" field.
func (sc *SecretCreate) SetAppID(i int64) *SecretCreate {
	sc.mutation.SetAppID(i)
	return sc
}

// SetActive sets the "active" field.
func (sc *SecretCreate) SetActive(b bool) *SecretCreate {
	sc.mutation.SetActive(b)
	return sc
}

// SetNillableActive sets the "active" field if the given value is not nil.
func (sc *SecretCreate) SetNillableActive(b *bool) *SecretCreate {
	if b != nil {
		sc.SetActive(*b)
	}
	return sc
}

// SetID sets the "id" field.
func (sc *SecretCreate) SetID(i int64) *SecretCreate {
	sc.mutation.SetID(i)
	return sc
}

// SetApp sets the "app" edge to the App entity.
func (sc *SecretCreate) SetApp(a *App) *SecretCreate {
	return sc.SetAppID(a.ID)
}

// Mutation returns the SecretMutation object of the builder.
func (sc *SecretCreate) Mutation() *SecretMutation {
	return sc.mutation
}

// Save creates the Secret in the database.
func (sc *SecretCreate) Save(ctx context.Context) (*Secret, error) {
	sc.defaults()
	return withHooks(ctx, sc.sqlSave, sc.mutation, sc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (sc *SecretCreate) SaveX(ctx context.Context) *Secret {
	v, err := sc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (sc *SecretCreate) Exec(ctx context.Context) error {
	_, err := sc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (sc *SecretCreate) ExecX(ctx context.Context) {
	if err := sc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (sc *SecretCreate) defaults() {
	if _, ok := sc.mutation.Active(); !ok {
		v := secret.DefaultActive
		sc.mutation.SetActive(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (sc *SecretCreate) check() error {
	if _, ok := sc.mutation.Key(); !ok {
		return &ValidationError{Name: "key", err: errors.New(`ent: missing required field "Secret.key"`)}
	}
	if _, ok := sc.mutation.Value(); !ok {
		return &ValidationError{Name: "value", err: errors.New(`ent: missing required field "Secret.value"`)}
	}
	if _, ok := sc.mutation.AppID(); !ok {
		return &ValidationError{Name: "app_id", err: errors.New(`ent: missing required field "Secret.app_id"`)}
	}
	if _, ok := sc.mutation.Active(); !ok {
		return &ValidationError{Name: "active", err: errors.New(`ent: missing required field "Secret.active"`)}
	}
	if len(sc.mutation.AppIDs()) == 0 {
		return &ValidationError{Name: "app", err: errors.New(`ent: missing required edge "Secret.app"`)}
	}
	return nil
}

func (sc *SecretCreate) sqlSave(ctx context.Context) (*Secret, error) {
	if err := sc.check(); err != nil {
		return nil, err
	}
	_node, _spec := sc.createSpec()
	if err := sqlgraph.CreateNode(ctx, sc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != _node.ID {
		id := _spec.ID.Value.(int64)
		_node.ID = int64(id)
	}
	sc.mutation.id = &_node.ID
	sc.mutation.done = true
	return _node, nil
}

func (sc *SecretCreate) createSpec() (*Secret, *sqlgraph.CreateSpec) {
	var (
		_node = &Secret{config: sc.config}
		_spec = sqlgraph.NewCreateSpec(secret.Table, sqlgraph.NewFieldSpec(secret.FieldID, field.TypeInt64))
	)
	if id, ok := sc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = id
	}
	if value, ok := sc.mutation.Key(); ok {
		_spec.SetField(secret.FieldKey, field.TypeString, value)
		_node.Key = value
	}
	if value, ok := sc.mutation.Value(); ok {
		_spec.SetField(secret.FieldValue, field.TypeString, value)
		_node.Value = value
	}
	if value, ok := sc.mutation.Active(); ok {
		_spec.SetField(secret.FieldActive, field.TypeBool, value)
		_node.Active = value
	}
	if nodes := sc.mutation.AppIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   secret.AppTable,
			Columns: []string{secret.AppColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(app.FieldID, field.TypeInt64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.AppID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// SecretCreateBulk is the builder for creating many Secret entities in bulk.
type SecretCreateBulk struct {
	config
	err      error
	builders []*SecretCreate
}

// Save creates the Secret entities in the database.
func (scb *SecretCreateBulk) Save(ctx context.Context) ([]*Secret, error) {
	if scb.err != nil {
		return nil, scb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(scb.builders))
	nodes := make([]*Secret, len(scb.builders))
	mutators := make([]Mutator, len(scb.builders))
	for i := range scb.builders {
		func(i int, root context.Context) {
			builder := scb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*SecretMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, scb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, scb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil && nodes[i].ID == 0 {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int64(id)
				}
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, scb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (scb *SecretCreateBulk) SaveX(ctx context.Context) []*Secret {
	v, err := scb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (scb *SecretCreateBulk) Exec(ctx context.Context) error {
	_, err := scb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (scb *SecretCreateBulk) ExecX(ctx context.Context) {
	if err := scb.Exec(ctx); err != nil {
		panic(err)
	}
}
