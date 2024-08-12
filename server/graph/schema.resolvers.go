package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"
	"fmt"

	"github.com/codecrafter404/bubble/graph/model"
)

// CreateOrder is the resolver for the createOrder field.
func (r *mutationResolver) CreateOrder(ctx context.Context, order model.NewOrder) (*model.Order, error) {
	panic(fmt.Errorf("not implemented: CreateOrder - createOrder"))
}

// UpdateOrder is the resolver for the updateOrder field.
func (r *mutationResolver) UpdateOrder(ctx context.Context, order int, state model.OrderState) (*model.Order, error) {
	panic(fmt.Errorf("not implemented: UpdateOrder - updateOrder"))
}

// DeleteOrder is the resolver for the deleteOrder field.
func (r *mutationResolver) DeleteOrder(ctx context.Context, order int) (int, error) {
	panic(fmt.Errorf("not implemented: DeleteOrder - deleteOrder"))
}

// UpdateItem is the resolver for the updateItem field.
func (r *mutationResolver) UpdateItem(ctx context.Context, id int, item model.UpdateItem) (*model.Item, error) {
	panic(fmt.Errorf("not implemented: UpdateItem - updateItem"))
}

// DeleteItems is the resolver for the deleteItems field.
func (r *mutationResolver) DeleteItems(ctx context.Context, id []int) ([]int, error) {
	panic(fmt.Errorf("not implemented: DeleteItems - deleteItems"))
}

// CreateItems is the resolver for the createItems field.
func (r *mutationResolver) CreateItems(ctx context.Context, items []*model.ItemInput) ([]int, error) {
	item_ids := []int{}
	for _, item := range items {
		item_ids = append(item_ids, item.ID)
	}

	rows, err := r.Db.Query("SELECT id FROM item WHERE id IN (?)", item_ids)
	if err != nil {
		return []int{}, fmt.Errorf("Failed to count rows: %w", err)
	}

	exItems := []int{}
	for rows.Next() {
		var exItemId int
		if err := rows.Scan(&exItemId); err != nil {
			return []int{}, fmt.Errorf("Failed to read row: %w", err)
		}
		exItems = append(exItems, exItemId)
	}
	tx, err := r.Db.Begin()
	if err != nil {
		return []int{}, fmt.Errorf("Failed to beginn transaction: %w", err)
	}
	for _, item := range items {
	}
	panic(fmt.Errorf("not implemented: CreateItems - createItems"))
}

// CreateCustomItems is the resolver for the createCustomItems field.
func (r *mutationResolver) CreateCustomItems(ctx context.Context, items []*model.CustomItemInput) ([]int, error) {
	panic(fmt.Errorf("not implemented: CreateCustomItems - createCustomItems"))
}

// GetPermission is the resolver for the getPermission field.
func (r *queryResolver) GetPermission(ctx context.Context) (model.User, error) {
	//TODO: implement authentication
	return model.UserUser, nil
}

// GetOrder is the resolver for the getOrder field.
func (r *queryResolver) GetOrder(ctx context.Context, id int) (*model.Order, error) {
	panic(fmt.Errorf("not implemented: GetOrder - getOrder"))
}

// Orders is the resolver for the orders field.
func (r *subscriptionResolver) Orders(ctx context.Context, state *model.OrderState, id *int, limit *int, skip *int) (<-chan []*model.Order, error) {
	panic(fmt.Errorf("not implemented: Orders - orders"))
}

// NextOrder is the resolver for the nextOrder field.
func (r *subscriptionResolver) NextOrder(ctx context.Context) (<-chan *model.Order, error) {
	panic(fmt.Errorf("not implemented: NextOrder - nextOrder"))
}

// Updates is the resolver for the updates field.
func (r *subscriptionResolver) Updates(ctx context.Context) (<-chan *model.UpdateEvent, error) {
	panic(fmt.Errorf("not implemented: Updates - updates"))
}

// Stats is the resolver for the stats field.
func (r *subscriptionResolver) Stats(ctx context.Context) (<-chan *model.Statistics, error) {
	panic(fmt.Errorf("not implemented: Stats - stats"))
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
