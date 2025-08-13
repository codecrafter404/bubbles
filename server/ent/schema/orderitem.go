package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// OrderItem holds the schema definition for the OrderItem entity.
type OrderItem struct {
	ent.Schema
}

// Fields of the OrderItem.
func (OrderItem) Fields() []ent.Field {
	return []ent.Field{
		field.Int("quantity"),
	}
}

// Edges of the OrderItem.
func (OrderItem) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("item", Item.Type).Unique(),
		edge.From("order", Order.Type).Ref("items").Unique(),
	}
}
