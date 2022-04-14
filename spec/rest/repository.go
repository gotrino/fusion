package rest

import (
	"context"
	"github.com/gotrino/fusion/runtime/rest"
	"github.com/gotrino/fusion/spec/app"
)

type Repository[T any] struct {
	Path    string // the resource path like /api/v1/books
	Default T
}

func (r Repository[T]) GetDefault() any {
	return r.Default
}

func (Repository[T]) IsRepository() bool {
	return true
}

func (r Repository[T]) New(ctx context.Context) app.RepositoryImplStencil {
	return rest.REST[T](ctx, r.Path).ToStencil()
}
