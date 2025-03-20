package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Profile struct {
	ent.Schema
}

func (Profile) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id"),
		field.String("firstname"),
		field.String("lastname"),
	}
}

func (Profile) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("profile").
			Unique(),
	}
}
