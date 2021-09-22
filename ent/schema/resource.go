package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Resource holds the schema definition for the Resource entity.
type Resource struct {
	ent.Schema
}

// Fields of the Resource.
func (Resource) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("path").Unique(),
		field.String("parent").Default("root"),
		field.Bool("router"),
		field.Bool("nav"),
		field.String("icon").Optional(),
		field.Uint("sort"),
	}
}

// Edges of the Resource.
func (Resource) Edges() []ent.Edge {
	return nil
}
