// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/errcode"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"helix.io/helix/ent/entemail"
	"helix.io/helix/ent/enttemporaryemail"
	"helix.io/helix/ent/entuser"
)

// Common entgql types.
type (
	Cursor         = entgql.Cursor[int]
	PageInfo       = entgql.PageInfo[int]
	OrderDirection = entgql.OrderDirection
)

func orderFunc(o OrderDirection, field string) func(*sql.Selector) {
	if o == entgql.OrderDirectionDesc {
		return Desc(field)
	}
	return Asc(field)
}

const errInvalidPagination = "INVALID_PAGINATION"

func validateFirstLast(first, last *int) (err *gqlerror.Error) {
	switch {
	case first != nil && last != nil:
		err = &gqlerror.Error{
			Message: "Passing both `first` and `last` to paginate a connection is not supported.",
		}
	case first != nil && *first < 0:
		err = &gqlerror.Error{
			Message: "`first` on a connection cannot be less than zero.",
		}
		errcode.Set(err, errInvalidPagination)
	case last != nil && *last < 0:
		err = &gqlerror.Error{
			Message: "`last` on a connection cannot be less than zero.",
		}
		errcode.Set(err, errInvalidPagination)
	}
	return err
}

func collectedField(ctx context.Context, path ...string) *graphql.CollectedField {
	fc := graphql.GetFieldContext(ctx)
	if fc == nil {
		return nil
	}
	field := fc.Field
	oc := graphql.GetOperationContext(ctx)
walk:
	for _, name := range path {
		for _, f := range graphql.CollectFields(oc, field.Selections, nil) {
			if f.Alias == name {
				field = f
				continue walk
			}
		}
		return nil
	}
	return &field
}

func hasCollectedField(ctx context.Context, path ...string) bool {
	if graphql.GetFieldContext(ctx) == nil {
		return true
	}
	return collectedField(ctx, path...) != nil
}

const (
	edgesField      = "edges"
	nodeField       = "node"
	pageInfoField   = "pageInfo"
	totalCountField = "totalCount"
)

func paginateLimit(first, last *int) int {
	var limit int
	if first != nil {
		limit = *first + 1
	} else if last != nil {
		limit = *last + 1
	}
	return limit
}

// EntEmailEdge is the edge representation of EntEmail.
type EntEmailEdge struct {
	Node   *EntEmail `json:"node"`
	Cursor Cursor    `json:"cursor"`
}

// EntEmailConnection is the connection containing edges to EntEmail.
type EntEmailConnection struct {
	Edges      []*EntEmailEdge `json:"edges"`
	PageInfo   PageInfo        `json:"pageInfo"`
	TotalCount int             `json:"totalCount"`
}

func (c *EntEmailConnection) build(nodes []*EntEmail, pager *entemailPager, after *Cursor, first *int, before *Cursor, last *int) {
	c.PageInfo.HasNextPage = before != nil
	c.PageInfo.HasPreviousPage = after != nil
	if first != nil && *first+1 == len(nodes) {
		c.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && *last+1 == len(nodes) {
		c.PageInfo.HasPreviousPage = true
		nodes = nodes[:len(nodes)-1]
	}
	var nodeAt func(int) *EntEmail
	if last != nil {
		n := len(nodes) - 1
		nodeAt = func(i int) *EntEmail {
			return nodes[n-i]
		}
	} else {
		nodeAt = func(i int) *EntEmail {
			return nodes[i]
		}
	}
	c.Edges = make([]*EntEmailEdge, len(nodes))
	for i := range nodes {
		node := nodeAt(i)
		c.Edges[i] = &EntEmailEdge{
			Node:   node,
			Cursor: pager.toCursor(node),
		}
	}
	if l := len(c.Edges); l > 0 {
		c.PageInfo.StartCursor = &c.Edges[0].Cursor
		c.PageInfo.EndCursor = &c.Edges[l-1].Cursor
	}
	if c.TotalCount == 0 {
		c.TotalCount = len(nodes)
	}
}

// EntEmailPaginateOption enables pagination customization.
type EntEmailPaginateOption func(*entemailPager) error

// WithEntEmailOrder configures pagination ordering.
func WithEntEmailOrder(order *EntEmailOrder) EntEmailPaginateOption {
	if order == nil {
		order = DefaultEntEmailOrder
	}
	o := *order
	return func(pager *entemailPager) error {
		if err := o.Direction.Validate(); err != nil {
			return err
		}
		if o.Field == nil {
			o.Field = DefaultEntEmailOrder.Field
		}
		pager.order = &o
		return nil
	}
}

// WithEntEmailFilter configures pagination filter.
func WithEntEmailFilter(filter func(*EntEmailQuery) (*EntEmailQuery, error)) EntEmailPaginateOption {
	return func(pager *entemailPager) error {
		if filter == nil {
			return errors.New("EntEmailQuery filter cannot be nil")
		}
		pager.filter = filter
		return nil
	}
}

type entemailPager struct {
	reverse bool
	order   *EntEmailOrder
	filter  func(*EntEmailQuery) (*EntEmailQuery, error)
}

func newEntEmailPager(opts []EntEmailPaginateOption, reverse bool) (*entemailPager, error) {
	pager := &entemailPager{reverse: reverse}
	for _, opt := range opts {
		if err := opt(pager); err != nil {
			return nil, err
		}
	}
	if pager.order == nil {
		pager.order = DefaultEntEmailOrder
	}
	return pager, nil
}

func (p *entemailPager) applyFilter(query *EntEmailQuery) (*EntEmailQuery, error) {
	if p.filter != nil {
		return p.filter(query)
	}
	return query, nil
}

func (p *entemailPager) toCursor(ee *EntEmail) Cursor {
	return p.order.Field.toCursor(ee)
}

func (p *entemailPager) applyCursors(query *EntEmailQuery, after, before *Cursor) (*EntEmailQuery, error) {
	direction := p.order.Direction
	if p.reverse {
		direction = direction.Reverse()
	}
	for _, predicate := range entgql.CursorsPredicate(after, before, DefaultEntEmailOrder.Field.column, p.order.Field.column, direction) {
		query = query.Where(predicate)
	}
	return query, nil
}

func (p *entemailPager) applyOrder(query *EntEmailQuery) *EntEmailQuery {
	direction := p.order.Direction
	if p.reverse {
		direction = direction.Reverse()
	}
	query = query.Order(p.order.Field.toTerm(direction.OrderTermOption()))
	if p.order.Field != DefaultEntEmailOrder.Field {
		query = query.Order(DefaultEntEmailOrder.Field.toTerm(direction.OrderTermOption()))
	}
	if len(query.ctx.Fields) > 0 {
		query.ctx.AppendFieldOnce(p.order.Field.column)
	}
	return query
}

func (p *entemailPager) orderExpr(query *EntEmailQuery) sql.Querier {
	direction := p.order.Direction
	if p.reverse {
		direction = direction.Reverse()
	}
	if len(query.ctx.Fields) > 0 {
		query.ctx.AppendFieldOnce(p.order.Field.column)
	}
	return sql.ExprFunc(func(b *sql.Builder) {
		b.Ident(p.order.Field.column).Pad().WriteString(string(direction))
		if p.order.Field != DefaultEntEmailOrder.Field {
			b.Comma().Ident(DefaultEntEmailOrder.Field.column).Pad().WriteString(string(direction))
		}
	})
}

// Paginate executes the query and returns a relay based cursor connection to EntEmail.
func (ee *EntEmailQuery) Paginate(
	ctx context.Context, after *Cursor, first *int,
	before *Cursor, last *int, opts ...EntEmailPaginateOption,
) (*EntEmailConnection, error) {
	if err := validateFirstLast(first, last); err != nil {
		return nil, err
	}
	pager, err := newEntEmailPager(opts, last != nil)
	if err != nil {
		return nil, err
	}
	if ee, err = pager.applyFilter(ee); err != nil {
		return nil, err
	}
	conn := &EntEmailConnection{Edges: []*EntEmailEdge{}}
	ignoredEdges := !hasCollectedField(ctx, edgesField)
	if hasCollectedField(ctx, totalCountField) || hasCollectedField(ctx, pageInfoField) {
		hasPagination := after != nil || first != nil || before != nil || last != nil
		if hasPagination || ignoredEdges {
			c := ee.Clone()
			c.ctx.Fields = nil
			if conn.TotalCount, err = c.Count(ctx); err != nil {
				return nil, err
			}
			conn.PageInfo.HasNextPage = first != nil && conn.TotalCount > 0
			conn.PageInfo.HasPreviousPage = last != nil && conn.TotalCount > 0
		}
	}
	if ignoredEdges || (first != nil && *first == 0) || (last != nil && *last == 0) {
		return conn, nil
	}
	if ee, err = pager.applyCursors(ee, after, before); err != nil {
		return nil, err
	}
	limit := paginateLimit(first, last)
	if limit != 0 {
		ee.Limit(limit)
	}
	if field := collectedField(ctx, edgesField, nodeField); field != nil {
		if err := ee.collectField(ctx, limit == 1, graphql.GetOperationContext(ctx), *field, []string{edgesField, nodeField}); err != nil {
			return nil, err
		}
	}
	ee = pager.applyOrder(ee)
	nodes, err := ee.All(ctx)
	if err != nil {
		return nil, err
	}
	conn.build(nodes, pager, after, first, before, last)
	return conn, nil
}

// EntEmailOrderField defines the ordering field of EntEmail.
type EntEmailOrderField struct {
	// Value extracts the ordering value from the given EntEmail.
	Value    func(*EntEmail) (ent.Value, error)
	column   string // field or computed.
	toTerm   func(...sql.OrderTermOption) entemail.OrderOption
	toCursor func(*EntEmail) Cursor
}

// EntEmailOrder defines the ordering of EntEmail.
type EntEmailOrder struct {
	Direction OrderDirection      `json:"direction"`
	Field     *EntEmailOrderField `json:"field"`
}

// DefaultEntEmailOrder is the default ordering of EntEmail.
var DefaultEntEmailOrder = &EntEmailOrder{
	Direction: entgql.OrderDirectionAsc,
	Field: &EntEmailOrderField{
		Value: func(ee *EntEmail) (ent.Value, error) {
			return ee.ID, nil
		},
		column: entemail.FieldID,
		toTerm: entemail.ByID,
		toCursor: func(ee *EntEmail) Cursor {
			return Cursor{ID: ee.ID}
		},
	},
}

// ToEdge converts EntEmail into EntEmailEdge.
func (ee *EntEmail) ToEdge(order *EntEmailOrder) *EntEmailEdge {
	if order == nil {
		order = DefaultEntEmailOrder
	}
	return &EntEmailEdge{
		Node:   ee,
		Cursor: order.Field.toCursor(ee),
	}
}

// EntTemporaryEmailEdge is the edge representation of EntTemporaryEmail.
type EntTemporaryEmailEdge struct {
	Node   *EntTemporaryEmail `json:"node"`
	Cursor Cursor             `json:"cursor"`
}

// EntTemporaryEmailConnection is the connection containing edges to EntTemporaryEmail.
type EntTemporaryEmailConnection struct {
	Edges      []*EntTemporaryEmailEdge `json:"edges"`
	PageInfo   PageInfo                 `json:"pageInfo"`
	TotalCount int                      `json:"totalCount"`
}

func (c *EntTemporaryEmailConnection) build(nodes []*EntTemporaryEmail, pager *enttemporaryemailPager, after *Cursor, first *int, before *Cursor, last *int) {
	c.PageInfo.HasNextPage = before != nil
	c.PageInfo.HasPreviousPage = after != nil
	if first != nil && *first+1 == len(nodes) {
		c.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && *last+1 == len(nodes) {
		c.PageInfo.HasPreviousPage = true
		nodes = nodes[:len(nodes)-1]
	}
	var nodeAt func(int) *EntTemporaryEmail
	if last != nil {
		n := len(nodes) - 1
		nodeAt = func(i int) *EntTemporaryEmail {
			return nodes[n-i]
		}
	} else {
		nodeAt = func(i int) *EntTemporaryEmail {
			return nodes[i]
		}
	}
	c.Edges = make([]*EntTemporaryEmailEdge, len(nodes))
	for i := range nodes {
		node := nodeAt(i)
		c.Edges[i] = &EntTemporaryEmailEdge{
			Node:   node,
			Cursor: pager.toCursor(node),
		}
	}
	if l := len(c.Edges); l > 0 {
		c.PageInfo.StartCursor = &c.Edges[0].Cursor
		c.PageInfo.EndCursor = &c.Edges[l-1].Cursor
	}
	if c.TotalCount == 0 {
		c.TotalCount = len(nodes)
	}
}

// EntTemporaryEmailPaginateOption enables pagination customization.
type EntTemporaryEmailPaginateOption func(*enttemporaryemailPager) error

// WithEntTemporaryEmailOrder configures pagination ordering.
func WithEntTemporaryEmailOrder(order *EntTemporaryEmailOrder) EntTemporaryEmailPaginateOption {
	if order == nil {
		order = DefaultEntTemporaryEmailOrder
	}
	o := *order
	return func(pager *enttemporaryemailPager) error {
		if err := o.Direction.Validate(); err != nil {
			return err
		}
		if o.Field == nil {
			o.Field = DefaultEntTemporaryEmailOrder.Field
		}
		pager.order = &o
		return nil
	}
}

// WithEntTemporaryEmailFilter configures pagination filter.
func WithEntTemporaryEmailFilter(filter func(*EntTemporaryEmailQuery) (*EntTemporaryEmailQuery, error)) EntTemporaryEmailPaginateOption {
	return func(pager *enttemporaryemailPager) error {
		if filter == nil {
			return errors.New("EntTemporaryEmailQuery filter cannot be nil")
		}
		pager.filter = filter
		return nil
	}
}

type enttemporaryemailPager struct {
	reverse bool
	order   *EntTemporaryEmailOrder
	filter  func(*EntTemporaryEmailQuery) (*EntTemporaryEmailQuery, error)
}

func newEntTemporaryEmailPager(opts []EntTemporaryEmailPaginateOption, reverse bool) (*enttemporaryemailPager, error) {
	pager := &enttemporaryemailPager{reverse: reverse}
	for _, opt := range opts {
		if err := opt(pager); err != nil {
			return nil, err
		}
	}
	if pager.order == nil {
		pager.order = DefaultEntTemporaryEmailOrder
	}
	return pager, nil
}

func (p *enttemporaryemailPager) applyFilter(query *EntTemporaryEmailQuery) (*EntTemporaryEmailQuery, error) {
	if p.filter != nil {
		return p.filter(query)
	}
	return query, nil
}

func (p *enttemporaryemailPager) toCursor(ete *EntTemporaryEmail) Cursor {
	return p.order.Field.toCursor(ete)
}

func (p *enttemporaryemailPager) applyCursors(query *EntTemporaryEmailQuery, after, before *Cursor) (*EntTemporaryEmailQuery, error) {
	direction := p.order.Direction
	if p.reverse {
		direction = direction.Reverse()
	}
	for _, predicate := range entgql.CursorsPredicate(after, before, DefaultEntTemporaryEmailOrder.Field.column, p.order.Field.column, direction) {
		query = query.Where(predicate)
	}
	return query, nil
}

func (p *enttemporaryemailPager) applyOrder(query *EntTemporaryEmailQuery) *EntTemporaryEmailQuery {
	direction := p.order.Direction
	if p.reverse {
		direction = direction.Reverse()
	}
	query = query.Order(p.order.Field.toTerm(direction.OrderTermOption()))
	if p.order.Field != DefaultEntTemporaryEmailOrder.Field {
		query = query.Order(DefaultEntTemporaryEmailOrder.Field.toTerm(direction.OrderTermOption()))
	}
	if len(query.ctx.Fields) > 0 {
		query.ctx.AppendFieldOnce(p.order.Field.column)
	}
	return query
}

func (p *enttemporaryemailPager) orderExpr(query *EntTemporaryEmailQuery) sql.Querier {
	direction := p.order.Direction
	if p.reverse {
		direction = direction.Reverse()
	}
	if len(query.ctx.Fields) > 0 {
		query.ctx.AppendFieldOnce(p.order.Field.column)
	}
	return sql.ExprFunc(func(b *sql.Builder) {
		b.Ident(p.order.Field.column).Pad().WriteString(string(direction))
		if p.order.Field != DefaultEntTemporaryEmailOrder.Field {
			b.Comma().Ident(DefaultEntTemporaryEmailOrder.Field.column).Pad().WriteString(string(direction))
		}
	})
}

// Paginate executes the query and returns a relay based cursor connection to EntTemporaryEmail.
func (ete *EntTemporaryEmailQuery) Paginate(
	ctx context.Context, after *Cursor, first *int,
	before *Cursor, last *int, opts ...EntTemporaryEmailPaginateOption,
) (*EntTemporaryEmailConnection, error) {
	if err := validateFirstLast(first, last); err != nil {
		return nil, err
	}
	pager, err := newEntTemporaryEmailPager(opts, last != nil)
	if err != nil {
		return nil, err
	}
	if ete, err = pager.applyFilter(ete); err != nil {
		return nil, err
	}
	conn := &EntTemporaryEmailConnection{Edges: []*EntTemporaryEmailEdge{}}
	ignoredEdges := !hasCollectedField(ctx, edgesField)
	if hasCollectedField(ctx, totalCountField) || hasCollectedField(ctx, pageInfoField) {
		hasPagination := after != nil || first != nil || before != nil || last != nil
		if hasPagination || ignoredEdges {
			c := ete.Clone()
			c.ctx.Fields = nil
			if conn.TotalCount, err = c.Count(ctx); err != nil {
				return nil, err
			}
			conn.PageInfo.HasNextPage = first != nil && conn.TotalCount > 0
			conn.PageInfo.HasPreviousPage = last != nil && conn.TotalCount > 0
		}
	}
	if ignoredEdges || (first != nil && *first == 0) || (last != nil && *last == 0) {
		return conn, nil
	}
	if ete, err = pager.applyCursors(ete, after, before); err != nil {
		return nil, err
	}
	limit := paginateLimit(first, last)
	if limit != 0 {
		ete.Limit(limit)
	}
	if field := collectedField(ctx, edgesField, nodeField); field != nil {
		if err := ete.collectField(ctx, limit == 1, graphql.GetOperationContext(ctx), *field, []string{edgesField, nodeField}); err != nil {
			return nil, err
		}
	}
	ete = pager.applyOrder(ete)
	nodes, err := ete.All(ctx)
	if err != nil {
		return nil, err
	}
	conn.build(nodes, pager, after, first, before, last)
	return conn, nil
}

// EntTemporaryEmailOrderField defines the ordering field of EntTemporaryEmail.
type EntTemporaryEmailOrderField struct {
	// Value extracts the ordering value from the given EntTemporaryEmail.
	Value    func(*EntTemporaryEmail) (ent.Value, error)
	column   string // field or computed.
	toTerm   func(...sql.OrderTermOption) enttemporaryemail.OrderOption
	toCursor func(*EntTemporaryEmail) Cursor
}

// EntTemporaryEmailOrder defines the ordering of EntTemporaryEmail.
type EntTemporaryEmailOrder struct {
	Direction OrderDirection               `json:"direction"`
	Field     *EntTemporaryEmailOrderField `json:"field"`
}

// DefaultEntTemporaryEmailOrder is the default ordering of EntTemporaryEmail.
var DefaultEntTemporaryEmailOrder = &EntTemporaryEmailOrder{
	Direction: entgql.OrderDirectionAsc,
	Field: &EntTemporaryEmailOrderField{
		Value: func(ete *EntTemporaryEmail) (ent.Value, error) {
			return ete.ID, nil
		},
		column: enttemporaryemail.FieldID,
		toTerm: enttemporaryemail.ByID,
		toCursor: func(ete *EntTemporaryEmail) Cursor {
			return Cursor{ID: ete.ID}
		},
	},
}

// ToEdge converts EntTemporaryEmail into EntTemporaryEmailEdge.
func (ete *EntTemporaryEmail) ToEdge(order *EntTemporaryEmailOrder) *EntTemporaryEmailEdge {
	if order == nil {
		order = DefaultEntTemporaryEmailOrder
	}
	return &EntTemporaryEmailEdge{
		Node:   ete,
		Cursor: order.Field.toCursor(ete),
	}
}

// EntUserEdge is the edge representation of EntUser.
type EntUserEdge struct {
	Node   *EntUser `json:"node"`
	Cursor Cursor   `json:"cursor"`
}

// EntUserConnection is the connection containing edges to EntUser.
type EntUserConnection struct {
	Edges      []*EntUserEdge `json:"edges"`
	PageInfo   PageInfo       `json:"pageInfo"`
	TotalCount int            `json:"totalCount"`
}

func (c *EntUserConnection) build(nodes []*EntUser, pager *entuserPager, after *Cursor, first *int, before *Cursor, last *int) {
	c.PageInfo.HasNextPage = before != nil
	c.PageInfo.HasPreviousPage = after != nil
	if first != nil && *first+1 == len(nodes) {
		c.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && *last+1 == len(nodes) {
		c.PageInfo.HasPreviousPage = true
		nodes = nodes[:len(nodes)-1]
	}
	var nodeAt func(int) *EntUser
	if last != nil {
		n := len(nodes) - 1
		nodeAt = func(i int) *EntUser {
			return nodes[n-i]
		}
	} else {
		nodeAt = func(i int) *EntUser {
			return nodes[i]
		}
	}
	c.Edges = make([]*EntUserEdge, len(nodes))
	for i := range nodes {
		node := nodeAt(i)
		c.Edges[i] = &EntUserEdge{
			Node:   node,
			Cursor: pager.toCursor(node),
		}
	}
	if l := len(c.Edges); l > 0 {
		c.PageInfo.StartCursor = &c.Edges[0].Cursor
		c.PageInfo.EndCursor = &c.Edges[l-1].Cursor
	}
	if c.TotalCount == 0 {
		c.TotalCount = len(nodes)
	}
}

// EntUserPaginateOption enables pagination customization.
type EntUserPaginateOption func(*entuserPager) error

// WithEntUserOrder configures pagination ordering.
func WithEntUserOrder(order *EntUserOrder) EntUserPaginateOption {
	if order == nil {
		order = DefaultEntUserOrder
	}
	o := *order
	return func(pager *entuserPager) error {
		if err := o.Direction.Validate(); err != nil {
			return err
		}
		if o.Field == nil {
			o.Field = DefaultEntUserOrder.Field
		}
		pager.order = &o
		return nil
	}
}

// WithEntUserFilter configures pagination filter.
func WithEntUserFilter(filter func(*EntUserQuery) (*EntUserQuery, error)) EntUserPaginateOption {
	return func(pager *entuserPager) error {
		if filter == nil {
			return errors.New("EntUserQuery filter cannot be nil")
		}
		pager.filter = filter
		return nil
	}
}

type entuserPager struct {
	reverse bool
	order   *EntUserOrder
	filter  func(*EntUserQuery) (*EntUserQuery, error)
}

func newEntUserPager(opts []EntUserPaginateOption, reverse bool) (*entuserPager, error) {
	pager := &entuserPager{reverse: reverse}
	for _, opt := range opts {
		if err := opt(pager); err != nil {
			return nil, err
		}
	}
	if pager.order == nil {
		pager.order = DefaultEntUserOrder
	}
	return pager, nil
}

func (p *entuserPager) applyFilter(query *EntUserQuery) (*EntUserQuery, error) {
	if p.filter != nil {
		return p.filter(query)
	}
	return query, nil
}

func (p *entuserPager) toCursor(eu *EntUser) Cursor {
	return p.order.Field.toCursor(eu)
}

func (p *entuserPager) applyCursors(query *EntUserQuery, after, before *Cursor) (*EntUserQuery, error) {
	direction := p.order.Direction
	if p.reverse {
		direction = direction.Reverse()
	}
	for _, predicate := range entgql.CursorsPredicate(after, before, DefaultEntUserOrder.Field.column, p.order.Field.column, direction) {
		query = query.Where(predicate)
	}
	return query, nil
}

func (p *entuserPager) applyOrder(query *EntUserQuery) *EntUserQuery {
	direction := p.order.Direction
	if p.reverse {
		direction = direction.Reverse()
	}
	query = query.Order(p.order.Field.toTerm(direction.OrderTermOption()))
	if p.order.Field != DefaultEntUserOrder.Field {
		query = query.Order(DefaultEntUserOrder.Field.toTerm(direction.OrderTermOption()))
	}
	if len(query.ctx.Fields) > 0 {
		query.ctx.AppendFieldOnce(p.order.Field.column)
	}
	return query
}

func (p *entuserPager) orderExpr(query *EntUserQuery) sql.Querier {
	direction := p.order.Direction
	if p.reverse {
		direction = direction.Reverse()
	}
	if len(query.ctx.Fields) > 0 {
		query.ctx.AppendFieldOnce(p.order.Field.column)
	}
	return sql.ExprFunc(func(b *sql.Builder) {
		b.Ident(p.order.Field.column).Pad().WriteString(string(direction))
		if p.order.Field != DefaultEntUserOrder.Field {
			b.Comma().Ident(DefaultEntUserOrder.Field.column).Pad().WriteString(string(direction))
		}
	})
}

// Paginate executes the query and returns a relay based cursor connection to EntUser.
func (eu *EntUserQuery) Paginate(
	ctx context.Context, after *Cursor, first *int,
	before *Cursor, last *int, opts ...EntUserPaginateOption,
) (*EntUserConnection, error) {
	if err := validateFirstLast(first, last); err != nil {
		return nil, err
	}
	pager, err := newEntUserPager(opts, last != nil)
	if err != nil {
		return nil, err
	}
	if eu, err = pager.applyFilter(eu); err != nil {
		return nil, err
	}
	conn := &EntUserConnection{Edges: []*EntUserEdge{}}
	ignoredEdges := !hasCollectedField(ctx, edgesField)
	if hasCollectedField(ctx, totalCountField) || hasCollectedField(ctx, pageInfoField) {
		hasPagination := after != nil || first != nil || before != nil || last != nil
		if hasPagination || ignoredEdges {
			c := eu.Clone()
			c.ctx.Fields = nil
			if conn.TotalCount, err = c.Count(ctx); err != nil {
				return nil, err
			}
			conn.PageInfo.HasNextPage = first != nil && conn.TotalCount > 0
			conn.PageInfo.HasPreviousPage = last != nil && conn.TotalCount > 0
		}
	}
	if ignoredEdges || (first != nil && *first == 0) || (last != nil && *last == 0) {
		return conn, nil
	}
	if eu, err = pager.applyCursors(eu, after, before); err != nil {
		return nil, err
	}
	limit := paginateLimit(first, last)
	if limit != 0 {
		eu.Limit(limit)
	}
	if field := collectedField(ctx, edgesField, nodeField); field != nil {
		if err := eu.collectField(ctx, limit == 1, graphql.GetOperationContext(ctx), *field, []string{edgesField, nodeField}); err != nil {
			return nil, err
		}
	}
	eu = pager.applyOrder(eu)
	nodes, err := eu.All(ctx)
	if err != nil {
		return nil, err
	}
	conn.build(nodes, pager, after, first, before, last)
	return conn, nil
}

// EntUserOrderField defines the ordering field of EntUser.
type EntUserOrderField struct {
	// Value extracts the ordering value from the given EntUser.
	Value    func(*EntUser) (ent.Value, error)
	column   string // field or computed.
	toTerm   func(...sql.OrderTermOption) entuser.OrderOption
	toCursor func(*EntUser) Cursor
}

// EntUserOrder defines the ordering of EntUser.
type EntUserOrder struct {
	Direction OrderDirection     `json:"direction"`
	Field     *EntUserOrderField `json:"field"`
}

// DefaultEntUserOrder is the default ordering of EntUser.
var DefaultEntUserOrder = &EntUserOrder{
	Direction: entgql.OrderDirectionAsc,
	Field: &EntUserOrderField{
		Value: func(eu *EntUser) (ent.Value, error) {
			return eu.ID, nil
		},
		column: entuser.FieldID,
		toTerm: entuser.ByID,
		toCursor: func(eu *EntUser) Cursor {
			return Cursor{ID: eu.ID}
		},
	},
}

// ToEdge converts EntUser into EntUserEdge.
func (eu *EntUser) ToEdge(order *EntUserOrder) *EntUserEdge {
	if order == nil {
		order = DefaultEntUserOrder
	}
	return &EntUserEdge{
		Node:   eu,
		Cursor: order.Field.toCursor(eu),
	}
}
