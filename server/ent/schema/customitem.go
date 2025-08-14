package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// CustomItem holds the schema definition for the CustomItem entity.
type CustomItem struct {
	ent.Schema
}

// Fields of the CustomItem.
func (CustomItem) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.Bool("exclusive").Comment("if true then one or more variants can be selected at once"),
		field.Int("prev").Immutable().Optional().Nillable().Comment("The id of the previous to selectable custom item"),
	}
}

// Edges of the CustomItem.
func (CustomItem) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("variants", Item.Type).Required(),
	}
}
func (CustomItem) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.QueryField(),
		entgql.Mutations(entgql.MutationCreate(), entgql.MutationUpdate()), entgql.QueryField(),
	}
}
