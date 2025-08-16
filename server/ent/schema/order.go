package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Order holds the schema definition for the Order entity.
type Order struct {
	ent.Schema
}

// Fields of the Order.
func (Order) Fields() []ent.Field {
	return []ent.Field{
		field.Time("submitted").Annotations(entgql.Skip(entgql.SkipMutationCreateInput, entgql.SkipMutationUpdateInput), entgql.OrderField("SUBMITTED")).Immutable().Comment("The time when the order has been submitted"),
		field.String("identifier").Immutable().Annotations(entgql.Skip(entgql.SkipMutationCreateInput, entgql.SkipMutationUpdateInput)).Comment("A sequencially generated string to identifiy an open order"),
		field.Enum("state").Values("created", "pending", "compleated", "cancelled").Default("created").Annotations(entgql.OrderField("STATE")),
		field.Float("total").Comment("on server generated orders total").Immutable().Annotations(entgql.Skip(entgql.SkipMutationCreateInput, entgql.SkipMutationUpdateInput), entgql.OrderField("TOTAL")),
	}
}

// Edges of the Order.
func (Order) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("items", OrderItem.Type).Immutable().Annotations(entgql.Skip(entgql.SkipMutationUpdateInput)),              //NOTE: This field SHOULD NOT be used as it will be overwritten
		edge.To("custom_items", OrderCustomItem.Type).Immutable().Annotations(entgql.Skip(entgql.SkipMutationUpdateInput)), //NOTE: This field SHOULD NOT be used as it will be overwritten
	}
}
func (Order) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.RelayConnection(),
		entgql.QueryField(),
		entgql.Mutations(entgql.MutationCreate(), entgql.MutationUpdate()),
	}
}
