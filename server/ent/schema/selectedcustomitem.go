package schema

import (
	"entgo.io/ent"
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
		edge.To("selected_variants", Item.Type),
		edge.To("custom_item", CustomItem.Type).Unique(),
	}
}
