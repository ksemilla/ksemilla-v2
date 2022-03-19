package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/ksemilla/ksemilla-v2/graph/generated"
	"github.com/ksemilla/ksemilla-v2/graph/model"
)

func (r *mutationResolver) CreateOwner(ctx context.Context, input model.CreateOwner) (*model.User, error) {
	return db.CreateOwner(&input)
}

func (r *mutationResolver) CreateUser(ctx context.Context, input model.CreateUser) (*model.User, error) {
	return db.CreateUser(&input)
}

func (r *mutationResolver) UpdateUser(ctx context.Context, input model.UpdateUser) (*model.User, error) {
	return db.UpdateUser(&input)
}

func (r *mutationResolver) DeleteUser(ctx context.Context, input string) (bool, error) {
	res, err := db.DeleteUser(input)
	return res.DeletedCount > 0, err
}

func (r *queryResolver) GetCreateOwnerValidity(ctx context.Context) (bool, error) {
	return db.IsCreateOwnerValid()
}

func (r *queryResolver) GetUsers(ctx context.Context) ([]*model.User, error) {
	return db.GetAllUsers()
}

func (r *queryResolver) FetchUser(ctx context.Context, id string) (*model.User, error) {
	return db.FindOneUser(id)
}

func (r *queryResolver) GetNumber(ctx context.Context) (int, error) {
	return 4, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
