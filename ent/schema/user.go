package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type User struct {
	ent.Schema
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id"),
		field.String("username"),
	}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("profile", Profile.Type).
			Unique(),
	}
}
