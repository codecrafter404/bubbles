package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// OrderCustomItem holds the schema definition for the OrderCustomItem entity.
type OrderCustomItem struct {
	ent.Schema
}

// Fields of the OrderCustomItem.
func (OrderCustomItem) Fields() []ent.Field {
	return []ent.Field{
		field.Int("quantity"),
	}
}

// Edges of the OrderCustomItem.
func (OrderCustomItem) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("selected_custom_items", SelectedCustomItem.Type),
		edge.To("master_custom_item", CustomItem.Type).Unique(),
		edge.From("order", Order.Type).Ref("custom_items").Unique(),
	}
}
