// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"helix.io/helix/ent/entemail"
	"helix.io/helix/ent/predicate"
)

// EntEmailDelete is the builder for deleting a EntEmail entity.
type EntEmailDelete struct {
	config
	hooks    []Hook
	mutation *EntEmailMutation
}

// Where appends a list predicates to the EntEmailDelete builder.
func (eed *EntEmailDelete) Where(ps ...predicate.EntEmail) *EntEmailDelete {
	eed.mutation.Where(ps...)
	return eed
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (eed *EntEmailDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, eed.sqlExec, eed.mutation, eed.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (eed *EntEmailDelete) ExecX(ctx context.Context) int {
	n, err := eed.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (eed *EntEmailDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(entemail.Table, sqlgraph.NewFieldSpec(entemail.FieldID, field.TypeInt))
	if ps := eed.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, eed.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	eed.mutation.done = true
	return affected, err
}

// EntEmailDeleteOne is the builder for deleting a single EntEmail entity.
type EntEmailDeleteOne struct {
	eed *EntEmailDelete
}

// Where appends a list predicates to the EntEmailDelete builder.
func (eedo *EntEmailDeleteOne) Where(ps ...predicate.EntEmail) *EntEmailDeleteOne {
	eedo.eed.mutation.Where(ps...)
	return eedo
}

// Exec executes the deletion query.
func (eedo *EntEmailDeleteOne) Exec(ctx context.Context) error {
	n, err := eedo.eed.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{entemail.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (eedo *EntEmailDeleteOne) ExecX(ctx context.Context) {
	if err := eedo.Exec(ctx); err != nil {
		panic(err)
	}
}
