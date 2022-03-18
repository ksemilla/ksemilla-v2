package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/ksemilla/ksemilla-v2/graph/model"
)

func (r *mutationResolver) CreateInvoice(ctx context.Context, input model.CreateInvoice) (*model.Invoice, error) {
	return db.CreateInvoice(&input), nil
}

func (r *mutationResolver) UpdateInvoice(ctx context.Context, input model.UpdateInvoice) (*model.Invoice, error) {
	return db.UpdateInvoice(&input)
}

func (r *mutationResolver) DeleteInvoice(ctx context.Context, id string) (bool, error) {
	res, err := db.DeleteInvoice(id)
	return res.DeletedCount > 0, err
}

func (r *queryResolver) GetInvoices(ctx context.Context, page int) (*model.PaginatedInvoices, error) {
	return db.GetAllInvoices(int64(page))
}

func (r *queryResolver) GetInvoice(ctx context.Context, id string) (*model.Invoice, error) {
	return db.GetInvoice(id)
}

func (r *queryResolver) Test(ctx context.Context) (bool, error) {
	db.Test()
	return false, nil
}
