// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/ent/app"
	"github.com/ent/apptype"
)

// AppTypeCreate is the builder for creating a AppType entity.
type AppTypeCreate struct {
	config
	mutation *AppTypeMutation
	hooks    []Hook
}

// SetName sets the "name" field.
func (atc *AppTypeCreate) SetName(s string) *AppTypeCreate {
	atc.mutation.SetName(s)
	return atc
}

// SetID sets the "id" field.
func (atc *AppTypeCreate) SetID(i int) *AppTypeCreate {
	atc.mutation.SetID(i)
	return atc
}

// AddAppIDs adds the "apps" edge to the App entity by IDs.
func (atc *AppTypeCreate) AddAppIDs(ids ...int64) *AppTypeCreate {
	atc.mutation.AddAppIDs(ids...)
	return atc
}

// AddApps adds the "apps" edges to the App entity.
func (atc *AppTypeCreate) AddApps(a ...*App) *AppTypeCreate {
	ids := make([]int64, len(a))
	for i := range a {
		ids[i] = a[i].ID
	}
	return atc.AddAppIDs(ids...)
}

// Mutation returns the AppTypeMutation object of the builder.
func (atc *AppTypeCreate) Mutation() *AppTypeMutation {
	return atc.mutation
}

// Save creates the AppType in the database.
func (atc *AppTypeCreate) Save(ctx context.Context) (*AppType, error) {
	var (
		err  error
		node *AppType
	)
	if len(atc.hooks) == 0 {
		if err = atc.check(); err != nil {
			return nil, err
		}
		node, err = atc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*AppTypeMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = atc.check(); err != nil {
				return nil, err
			}
			atc.mutation = mutation
			if node, err = atc.sqlSave(ctx); err != nil {
				return nil, err
			}
			mutation.id = &node.ID
			mutation.done = true
			return node, err
		})
		for i := len(atc.hooks) - 1; i >= 0; i-- {
			if atc.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = atc.hooks[i](mut)
		}
		v, err := mut.Mutate(ctx, atc.mutation)
		if err != nil {
			return nil, err
		}
		nv, ok := v.(*AppType)
		if !ok {
			return nil, fmt.Errorf("unexpected node type %T returned from AppTypeMutation", v)
		}
		node = nv
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (atc *AppTypeCreate) SaveX(ctx context.Context) *AppType {
	v, err := atc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (atc *AppTypeCreate) Exec(ctx context.Context) error {
	_, err := atc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (atc *AppTypeCreate) ExecX(ctx context.Context) {
	if err := atc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (atc *AppTypeCreate) check() error {
	if _, ok := atc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`ent: missing required field "AppTypeID.name"`)}
	}
	return nil
}

func (atc *AppTypeCreate) sqlSave(ctx context.Context) (*AppType, error) {
	_node, _spec := atc.createSpec()
	if err := sqlgraph.CreateNode(ctx, atc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != _node.ID {
		id := _spec.ID.Value.(int64)
		_node.ID = int(id)
	}
	return _node, nil
}

func (atc *AppTypeCreate) createSpec() (*AppType, *sqlgraph.CreateSpec) {
	var (
		_node = &AppType{config: atc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: apptype.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: apptype.FieldID,
			},
		}
	)
	if id, ok := atc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = id
	}
	if value, ok := atc.mutation.Name(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: apptype.FieldName,
		})
		_node.Name = value
	}
	if nodes := atc.mutation.AppsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   apptype.AppsTable,
			Columns: []string{apptype.AppsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt64,
					Column: app.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// AppTypeCreateBulk is the builder for creating many AppType entities in bulk.
type AppTypeCreateBulk struct {
	config
	builders []*AppTypeCreate
}

// Save creates the AppType entities in the database.
func (atcb *AppTypeCreateBulk) Save(ctx context.Context) ([]*AppType, error) {
	specs := make([]*sqlgraph.CreateSpec, len(atcb.builders))
	nodes := make([]*AppType, len(atcb.builders))
	mutators := make([]Mutator, len(atcb.builders))
	for i := range atcb.builders {
		func(i int, root context.Context) {
			builder := atcb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*AppTypeMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				nodes[i], specs[i] = builder.createSpec()
				var err error
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, atcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, atcb.driver, spec); err != nil {
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
					nodes[i].ID = int(id)
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
		if _, err := mutators[0].Mutate(ctx, atcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (atcb *AppTypeCreateBulk) SaveX(ctx context.Context) []*AppType {
	v, err := atcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (atcb *AppTypeCreateBulk) Exec(ctx context.Context) error {
	_, err := atcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (atcb *AppTypeCreateBulk) ExecX(ctx context.Context) {
	if err := atcb.Exec(ctx); err != nil {
		panic(err)
	}
}
