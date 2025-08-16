package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
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
		edge.To("selected_custom_items", SelectedCustomItem.Type).Immutable(), //NOTE: This field SHOULD NOT be used as it will be overwritten
		edge.To("master_custom_item", CustomItem.Type).Unique().Immutable().Required(),
		edge.From("order", Order.Type).Ref("custom_items").Unique().Annotations(entgql.Skip(entgql.SkipMutationCreateInput)),
	}
}
func (OrderCustomItem) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.Mutations(entgql.MutationCreate()),
	}
}
