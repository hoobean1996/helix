// Code generated by ent, DO NOT EDIT.

package enttemporaryemail

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"helix.io/helix/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.EntTemporaryEmail {
	return predicate.EntTemporaryEmail(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.EntTemporaryEmail {
	return predicate.EntTemporaryEmail(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.EntTemporaryEmail {
	return predicate.EntTemporaryEmail(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.EntTemporaryEmail {
	return predicate.EntTemporaryEmail(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.EntTemporaryEmail {
	return predicate.EntTemporaryEmail(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.EntTemporaryEmail {
	return predicate.EntTemporaryEmail(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.EntTemporaryEmail {
	return predicate.EntTemporaryEmail(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.EntTemporaryEmail {
	return predicate.EntTemporaryEmail(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.EntTemporaryEmail {
	return predicate.EntTemporaryEmail(sql.FieldLTE(FieldID, id))
}

// Email applies equality check predicate on the "email" field. It's identical to EmailEQ.
func Email(v string) predicate.EntTemporaryEmail {
	return predicate.EntTemporaryEmail(sql.FieldEQ(FieldEmail, v))
}

// EmailEQ applies the EQ predicate on the "email" field.
func EmailEQ(v string) predicate.EntTemporaryEmail {
	return predicate.EntTemporaryEmail(sql.FieldEQ(FieldEmail, v))
}

// EmailNEQ applies the NEQ predicate on the "email" field.
func EmailNEQ(v string) predicate.EntTemporaryEmail {
	return predicate.EntTemporaryEmail(sql.FieldNEQ(FieldEmail, v))
}

// EmailIn applies the In predicate on the "email" field.
func EmailIn(vs ...string) predicate.EntTemporaryEmail {
	return predicate.EntTemporaryEmail(sql.FieldIn(FieldEmail, vs...))
}

// EmailNotIn applies the NotIn predicate on the "email" field.
func EmailNotIn(vs ...string) predicate.EntTemporaryEmail {
	return predicate.EntTemporaryEmail(sql.FieldNotIn(FieldEmail, vs...))
}

// EmailGT applies the GT predicate on the "email" field.
func EmailGT(v string) predicate.EntTemporaryEmail {
	return predicate.EntTemporaryEmail(sql.FieldGT(FieldEmail, v))
}

// EmailGTE applies the GTE predicate on the "email" field.
func EmailGTE(v string) predicate.EntTemporaryEmail {
	return predicate.EntTemporaryEmail(sql.FieldGTE(FieldEmail, v))
}

// EmailLT applies the LT predicate on the "email" field.
func EmailLT(v string) predicate.EntTemporaryEmail {
	return predicate.EntTemporaryEmail(sql.FieldLT(FieldEmail, v))
}

// EmailLTE applies the LTE predicate on the "email" field.
func EmailLTE(v string) predicate.EntTemporaryEmail {
	return predicate.EntTemporaryEmail(sql.FieldLTE(FieldEmail, v))
}

// EmailContains applies the Contains predicate on the "email" field.
func EmailContains(v string) predicate.EntTemporaryEmail {
	return predicate.EntTemporaryEmail(sql.FieldContains(FieldEmail, v))
}

// EmailHasPrefix applies the HasPrefix predicate on the "email" field.
func EmailHasPrefix(v string) predicate.EntTemporaryEmail {
	return predicate.EntTemporaryEmail(sql.FieldHasPrefix(FieldEmail, v))
}

// EmailHasSuffix applies the HasSuffix predicate on the "email" field.
func EmailHasSuffix(v string) predicate.EntTemporaryEmail {
	return predicate.EntTemporaryEmail(sql.FieldHasSuffix(FieldEmail, v))
}

// EmailEqualFold applies the EqualFold predicate on the "email" field.
func EmailEqualFold(v string) predicate.EntTemporaryEmail {
	return predicate.EntTemporaryEmail(sql.FieldEqualFold(FieldEmail, v))
}

// EmailContainsFold applies the ContainsFold predicate on the "email" field.
func EmailContainsFold(v string) predicate.EntTemporaryEmail {
	return predicate.EntTemporaryEmail(sql.FieldContainsFold(FieldEmail, v))
}

// HasOwner applies the HasEdge predicate on the "owner" edge.
func HasOwner() predicate.EntTemporaryEmail {
	return predicate.EntTemporaryEmail(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, OwnerTable, OwnerColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasOwnerWith applies the HasEdge predicate on the "owner" edge with a given conditions (other predicates).
func HasOwnerWith(preds ...predicate.EntUser) predicate.EntTemporaryEmail {
	return predicate.EntTemporaryEmail(func(s *sql.Selector) {
		step := newOwnerStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.EntTemporaryEmail) predicate.EntTemporaryEmail {
	return predicate.EntTemporaryEmail(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.EntTemporaryEmail) predicate.EntTemporaryEmail {
	return predicate.EntTemporaryEmail(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.EntTemporaryEmail) predicate.EntTemporaryEmail {
	return predicate.EntTemporaryEmail(sql.NotPredicates(p))
}
