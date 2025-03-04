package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// EntUser holds the schema definition for the EntUser entity.
type EntUser struct {
	ent.Schema
}

// Fields of the EntUser.
func (EntUser) Fields() []ent.Field {
	return []ent.Field{
		field.String("account").Comment("user's real email"),
	}
}

// Edges of the EntUser.
func (EntUser) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("owns", EntTemporaryEmail.Type),
	}
}

func (EntUser) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.QueryField(),
		entgql.Mutations(entgql.MutationCreate()),
	}
}
