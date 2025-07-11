package resolvers

// THIS CODE WILL BE UPDATED WITH SCHEMA CHANGES. PREVIOUS IMPLEMENTATION FOR SCHEMA CHANGES WILL BE KEPT IN THE COMMENT SECTION. IMPLEMENTATION FOR UNCHANGED SCHEMA WILL BE KEPT.

import (
	"context"

	"github.com/codecrafter404/bubble/graph"
	"github.com/codecrafter404/bubble/resolvers/mutation"
	"github.com/codecrafter404/bubble/resolvers/query"
	"github.com/codecrafter404/bubble/resolvers/subscription"
	"gorm.io/gorm"
)

type Resolver struct {
	EventChannel []chan *graph.UpdateEvent
	Db           *gorm.DB
}

// CreateOrder is the resolver for the createOrder field.
func (r *mutationResolver) CreateOrder(ctx context.Context, order graph.NewOrder) (*graph.Order, error) {
	return mutation.CreateOrder(ctx, order)
}

// UpdateOrder is the resolver for the updateOrder field.
func (r *mutationResolver) UpdateOrder(ctx context.Context, order int, state graph.OrderState) (*graph.Order, error) {
	return mutation.UpdateOrder(ctx, order, state)
}

// DeleteOrder is the resolver for the deleteOrder field.
func (r *mutationResolver) DeleteOrder(ctx context.Context, order int) (int, error) {
	return mutation.DeleteOrder(ctx, order)
}

// UpdateItem is the resolver for the updateItem field.
func (r *mutationResolver) UpdateItem(ctx context.Context, id int, item graph.UpdateItem) (*graph.Item, error) {
	return mutation.UpdateItem(ctx, id, item)
}

// UpdateCustomItem is the resolver for the updateCustomItem field.
func (r *mutationResolver) UpdateCustomItem(ctx context.Context, id int, item graph.UpdateCustomItem) (*graph.CustomItem, error) {
	return mutation.UpdateCustomItem(ctx, id, item)
}

// CreateItems is the resolver for the createItems field.
func (r *mutationResolver) CreateItems(ctx context.Context, items []*graph.NewItem) ([]int, error) {
	return mutation.CreateItems(ctx, items)
}

// CreateCustomItems is the resolver for the createCustomItems field.
func (r *mutationResolver) CreateCustomItems(ctx context.Context, items []*graph.NewCustomItem) ([]int, error) {
	return mutation.CreateCustomItems(ctx, items)
}

// GetPermission is the resolver for the getPermission field.
func (r *queryResolver) GetPermission(ctx context.Context) (graph.User, error) {
	return query.GetPermission(ctx)
}

// GetOrder is the resolver for the getOrder field.
func (r *queryResolver) GetOrder(ctx context.Context, id int) (*graph.Order, error) {
	return query.GetOrder(ctx, id)
}

// GetItems is the resolver for the getItems field.
func (r *queryResolver) GetItems(ctx context.Context) ([]*graph.Item, error) {
	return query.GetItems(ctx)
}

// GetCustomItems is the resolver for the getCustomItems field.
func (r *queryResolver) GetCustomItems(ctx context.Context) ([]*graph.CustomItem, error) {
	return query.GetCustomItems(ctx)
}

// Orders is the resolver for the orders field.
func (r *subscriptionResolver) Orders(ctx context.Context, state *graph.OrderState, id *int, limit *int, skip *int, sortAsc *bool) (<-chan []*graph.Order, error) {
	return subscription.Orders(ctx, state, id, limit, skip, sortAsc)
}

// NextOrder is the resolver for the nextOrder field.
func (r *subscriptionResolver) NextOrder(ctx context.Context) (<-chan *graph.Order, error) {
	return subscription.NextOrder(ctx)
}

// Updates is the resolver for the updates field.
func (r *subscriptionResolver) Updates(ctx context.Context) (<-chan *graph.UpdateEvent, error) {
	return subscription.Updates(ctx)
}

// Stats is the resolver for the stats field.
func (r *subscriptionResolver) Stats(ctx context.Context) (<-chan *graph.Statistics, error) {
	return subscription.Stats(ctx)
}

// Mutation returns graph.MutationResolver implementation.
func (r *Resolver) Mutation() graph.MutationResolver { return &mutationResolver{r} }

// Query returns graph.QueryResolver implementation.
func (r *Resolver) Query() graph.QueryResolver { return &queryResolver{r} }

// Subscription returns graph.SubscriptionResolver implementation.
func (r *Resolver) Subscription() graph.SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
/*
	type Resolver struct{}
*/
