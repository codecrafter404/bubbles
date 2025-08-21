package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Item holds the schema definition for the Item entity.
type Item struct {
	ent.Schema
}

// Fields of the Item.
func (Item) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Comment("eg.: cheesecake"),
		field.Float("price"),
		field.String("image").Comment("an url; eg.: https://example.com/cheesecake.png"),
		field.Bool("in_stock").Comment("weather the item is (physically) in stock"),
		field.String("notes").Comment("Custom text, which will be displayed"),
	}
}

// Edges of the Item.
func (Item) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("selected_custom_items", SelectedCustomItem.Type).Ref("selected_variants").Annotations(entgql.Skip(entgql.SkipAll)),
	}
}
func (Item) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.QueryField(),
		entgql.Mutations(entgql.MutationCreate(), entgql.MutationUpdate()), entgql.QueryField(),
	}
}
