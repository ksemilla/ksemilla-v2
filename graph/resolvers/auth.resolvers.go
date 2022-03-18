package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/ksemilla/ksemilla-v2/database"
	"github.com/ksemilla/ksemilla-v2/graph/model"
)

func (r *mutationResolver) Login(ctx context.Context, input model.LoginInput) (*model.LoginReturn, error) {
	return db.Login(&input)
}

func (r *mutationResolver) VerifyToken(ctx context.Context, input model.VerifyToken) (*model.User, error) {
	return db.VerifyToken(&input)
}

func (r *mutationResolver) ChangePassword(ctx context.Context, input model.ChangePassword) (bool, error) {
	return db.ChangePassword(&input)
}

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
var db = database.Connect()
