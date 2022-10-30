// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"database/sql/driver"
	"fmt"
	"math"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/ent/app"
	"github.com/ent/apptype"
	"github.com/ent/predicate"
)

// AppTypeQuery is the builder for querying AppType entities.
type AppTypeQuery struct {
	config
	limit      *int
	offset     *int
	unique     *bool
	order      []OrderFunc
	fields     []string
	predicates []predicate.AppType
	withApps   *AppQuery
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the AppTypeQuery builder.
func (atq *AppTypeQuery) Where(ps ...predicate.AppType) *AppTypeQuery {
	atq.predicates = append(atq.predicates, ps...)
	return atq
}

// Limit adds a limit step to the query.
func (atq *AppTypeQuery) Limit(limit int) *AppTypeQuery {
	atq.limit = &limit
	return atq
}

// Offset adds an offset step to the query.
func (atq *AppTypeQuery) Offset(offset int) *AppTypeQuery {
	atq.offset = &offset
	return atq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (atq *AppTypeQuery) Unique(unique bool) *AppTypeQuery {
	atq.unique = &unique
	return atq
}

// Order adds an order step to the query.
func (atq *AppTypeQuery) Order(o ...OrderFunc) *AppTypeQuery {
	atq.order = append(atq.order, o...)
	return atq
}

// QueryApps chains the current query on the "apps" edge.
func (atq *AppTypeQuery) QueryApps() *AppQuery {
	query := &AppQuery{config: atq.config}
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := atq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := atq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(apptype.Table, apptype.FieldID, selector),
			sqlgraph.To(app.Table, app.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, apptype.AppsTable, apptype.AppsColumn),
		)
		fromU = sqlgraph.SetNeighbors(atq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first AppType entity from the query.
// Returns a *NotFoundError when no AppType was found.
func (atq *AppTypeQuery) First(ctx context.Context) (*AppType, error) {
	nodes, err := atq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{apptype.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (atq *AppTypeQuery) FirstX(ctx context.Context) *AppType {
	node, err := atq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first AppType ID from the query.
// Returns a *NotFoundError when no AppType ID was found.
func (atq *AppTypeQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = atq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{apptype.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (atq *AppTypeQuery) FirstIDX(ctx context.Context) int {
	id, err := atq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single AppType entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one AppType entity is found.
// Returns a *NotFoundError when no AppType entities are found.
func (atq *AppTypeQuery) Only(ctx context.Context) (*AppType, error) {
	nodes, err := atq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{apptype.Label}
	default:
		return nil, &NotSingularError{apptype.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (atq *AppTypeQuery) OnlyX(ctx context.Context) *AppType {
	node, err := atq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only AppType ID in the query.
// Returns a *NotSingularError when more than one AppType ID is found.
// Returns a *NotFoundError when no entities are found.
func (atq *AppTypeQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = atq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{apptype.Label}
	default:
		err = &NotSingularError{apptype.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (atq *AppTypeQuery) OnlyIDX(ctx context.Context) int {
	id, err := atq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of AppTypes.
func (atq *AppTypeQuery) All(ctx context.Context) ([]*AppType, error) {
	if err := atq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return atq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (atq *AppTypeQuery) AllX(ctx context.Context) []*AppType {
	nodes, err := atq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of AppType IDs.
func (atq *AppTypeQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := atq.Select(apptype.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (atq *AppTypeQuery) IDsX(ctx context.Context) []int {
	ids, err := atq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (atq *AppTypeQuery) Count(ctx context.Context) (int, error) {
	if err := atq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return atq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (atq *AppTypeQuery) CountX(ctx context.Context) int {
	count, err := atq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (atq *AppTypeQuery) Exist(ctx context.Context) (bool, error) {
	if err := atq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return atq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (atq *AppTypeQuery) ExistX(ctx context.Context) bool {
	exist, err := atq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the AppTypeQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (atq *AppTypeQuery) Clone() *AppTypeQuery {
	if atq == nil {
		return nil
	}
	return &AppTypeQuery{
		config:     atq.config,
		limit:      atq.limit,
		offset:     atq.offset,
		order:      append([]OrderFunc{}, atq.order...),
		predicates: append([]predicate.AppType{}, atq.predicates...),
		withApps:   atq.withApps.Clone(),
		// clone intermediate query.
		sql:    atq.sql.Clone(),
		path:   atq.path,
		unique: atq.unique,
	}
}

// WithApps tells the query-builder to eager-load the nodes that are connected to
// the "apps" edge. The optional arguments are used to configure the query builder of the edge.
func (atq *AppTypeQuery) WithApps(opts ...func(*AppQuery)) *AppTypeQuery {
	query := &AppQuery{config: atq.config}
	for _, opt := range opts {
		opt(query)
	}
	atq.withApps = query
	return atq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Name string `json:"name,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.AppTypeID.Query().
//		GroupBy(apptype.FieldName).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (atq *AppTypeQuery) GroupBy(field string, fields ...string) *AppTypeGroupBy {
	grbuild := &AppTypeGroupBy{config: atq.config}
	grbuild.fields = append([]string{field}, fields...)
	grbuild.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := atq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return atq.sqlQuery(ctx), nil
	}
	grbuild.label = apptype.Label
	grbuild.flds, grbuild.scan = &grbuild.fields, grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		Name string `json:"name,omitempty"`
//	}
//
//	client.AppTypeID.Query().
//		Select(apptype.FieldName).
//		Scan(ctx, &v)
func (atq *AppTypeQuery) Select(fields ...string) *AppTypeSelect {
	atq.fields = append(atq.fields, fields...)
	selbuild := &AppTypeSelect{AppTypeQuery: atq}
	selbuild.label = apptype.Label
	selbuild.flds, selbuild.scan = &atq.fields, selbuild.Scan
	return selbuild
}

func (atq *AppTypeQuery) prepareQuery(ctx context.Context) error {
	for _, f := range atq.fields {
		if !apptype.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if atq.path != nil {
		prev, err := atq.path(ctx)
		if err != nil {
			return err
		}
		atq.sql = prev
	}
	return nil
}

func (atq *AppTypeQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*AppType, error) {
	var (
		nodes       = []*AppType{}
		_spec       = atq.querySpec()
		loadedTypes = [1]bool{
			atq.withApps != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*AppType).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &AppType{config: atq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, atq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := atq.withApps; query != nil {
		if err := atq.loadApps(ctx, query, nodes,
			func(n *AppType) { n.Edges.Apps = []*App{} },
			func(n *AppType, e *App) { n.Edges.Apps = append(n.Edges.Apps, e) }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (atq *AppTypeQuery) loadApps(ctx context.Context, query *AppQuery, nodes []*AppType, init func(*AppType), assign func(*AppType, *App)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[int]*AppType)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	query.Where(predicate.App(func(s *sql.Selector) {
		s.Where(sql.InValues(apptype.AppsColumn, fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.AppTypeID
		node, ok := nodeids[fk]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "app_type_id" returned %v for node %v`, fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}

func (atq *AppTypeQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := atq.querySpec()
	_spec.Node.Columns = atq.fields
	if len(atq.fields) > 0 {
		_spec.Unique = atq.unique != nil && *atq.unique
	}
	return sqlgraph.CountNodes(ctx, atq.driver, _spec)
}

func (atq *AppTypeQuery) sqlExist(ctx context.Context) (bool, error) {
	switch _, err := atq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

func (atq *AppTypeQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   apptype.Table,
			Columns: apptype.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: apptype.FieldID,
			},
		},
		From:   atq.sql,
		Unique: true,
	}
	if unique := atq.unique; unique != nil {
		_spec.Unique = *unique
	}
	if fields := atq.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, apptype.FieldID)
		for i := range fields {
			if fields[i] != apptype.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := atq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := atq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := atq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := atq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (atq *AppTypeQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(atq.driver.Dialect())
	t1 := builder.Table(apptype.Table)
	columns := atq.fields
	if len(columns) == 0 {
		columns = apptype.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if atq.sql != nil {
		selector = atq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if atq.unique != nil && *atq.unique {
		selector.Distinct()
	}
	for _, p := range atq.predicates {
		p(selector)
	}
	for _, p := range atq.order {
		p(selector)
	}
	if offset := atq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := atq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// AppTypeGroupBy is the group-by builder for AppType entities.
type AppTypeGroupBy struct {
	config
	selector
	fields []string
	fns    []AggregateFunc
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (atgb *AppTypeGroupBy) Aggregate(fns ...AggregateFunc) *AppTypeGroupBy {
	atgb.fns = append(atgb.fns, fns...)
	return atgb
}

// Scan applies the group-by query and scans the result into the given value.
func (atgb *AppTypeGroupBy) Scan(ctx context.Context, v any) error {
	query, err := atgb.path(ctx)
	if err != nil {
		return err
	}
	atgb.sql = query
	return atgb.sqlScan(ctx, v)
}

func (atgb *AppTypeGroupBy) sqlScan(ctx context.Context, v any) error {
	for _, f := range atgb.fields {
		if !apptype.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("invalid field %q for group-by", f)}
		}
	}
	selector := atgb.sqlQuery()
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := atgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (atgb *AppTypeGroupBy) sqlQuery() *sql.Selector {
	selector := atgb.sql.Select()
	aggregation := make([]string, 0, len(atgb.fns))
	for _, fn := range atgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	// If no columns were selected in a custom aggregation function, the default
	// selection is the fields used for "group-by", and the aggregation functions.
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(atgb.fields)+len(atgb.fns))
		for _, f := range atgb.fields {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	return selector.GroupBy(selector.Columns(atgb.fields...)...)
}

// AppTypeSelect is the builder for selecting fields of AppType entities.
type AppTypeSelect struct {
	*AppTypeQuery
	selector
	// intermediate query (i.e. traversal path).
	sql *sql.Selector
}

// Scan applies the selector query and scans the result into the given value.
func (ats *AppTypeSelect) Scan(ctx context.Context, v any) error {
	if err := ats.prepareQuery(ctx); err != nil {
		return err
	}
	ats.sql = ats.AppTypeQuery.sqlQuery(ctx)
	return ats.sqlScan(ctx, v)
}

func (ats *AppTypeSelect) sqlScan(ctx context.Context, v any) error {
	rows := &sql.Rows{}
	query, args := ats.sql.Query()
	if err := ats.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
