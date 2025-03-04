package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// EntEmail holds the schema definition for the EntEmail entity.
type EntEmail struct {
	ent.Schema
}

// Fields of the EntEmail.
func (EntEmail) Fields() []ent.Field {
	return []ent.Field{
		field.String("from"),
		field.Strings("to"),
		field.String("date"),
		field.String("topic"),
		field.String("body"),
	}
}

// Edges of the EntEmail.
func (EntEmail) Edges() []ent.Edge {
	return nil
}

func (EntEmail) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.QueryField(),
		entgql.Mutations(entgql.MutationCreate()),
	}
}
