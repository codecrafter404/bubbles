package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
)

// SelectedCustomItem holds the schema definition for the SelectedCustomItem entity.
type SelectedCustomItem struct {
	ent.Schema
}

// Fields of the SelectedCustomItem.
func (SelectedCustomItem) Fields() []ent.Field {
	return nil
}

// Edges of the SelectedCustomItem.
func (SelectedCustomItem) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("selected_variants", Item.Type).Immutable().Required(),
		edge.To("custom_item", CustomItem.Type).Unique().Immutable().Required(),
	}
}
func (SelectedCustomItem) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.Mutations(entgql.MutationCreate()),
	}
}
