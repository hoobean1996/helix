package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// EntTemporaryEmail holds the schema definition for the EntTemporaryEmail entity.
type EntTemporaryEmail struct {
	ent.Schema
}

// Fields of the EntTemporaryEmail.
func (EntTemporaryEmail) Fields() []ent.Field {
	return []ent.Field{
		field.String("email").Comment("test@vmail.today"),
	}
}

// Edges of the EntTemporaryEmail.
func (EntTemporaryEmail) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", EntUser.Type).Ref("owns").Unique(),
	}
}

func (EntTemporaryEmail) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.QueryField(),
		entgql.Mutations(entgql.MutationCreate()),
	}
}
