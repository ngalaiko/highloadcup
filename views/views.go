package views

import (
	"context"

	"github.com/ngalayko/highloadcup/database"
)

const (
	viewsCtxKey ctxKey = "ctx_key_for_views"
)

type ctxKey string

type Views struct {
	db *database.DB
}

// NewContext stores views in context
func NewContext(ctx context.Context, views interface{}) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	if views == nil {
		views = NewViews(ctx)
	}

	return context.WithValue(ctx, viewsCtxKey, views)
}

// FromContext returns views from context
func FromContext(ctx context.Context) *Views {
	if views, ok := ctx.Value(viewsCtxKey).(*Views); ok {
		return views
	}

	return NewViews(ctx)
}

// NewViews is a views constructor
func NewViews(ctx context.Context) *Views {
	return &Views{
		db: database.FromContext(ctx),
	}
}
