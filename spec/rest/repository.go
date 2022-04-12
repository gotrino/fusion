package rest

import (
	"context"
	"github.com/gotrino/fusion/runtime/rest"
	"github.com/gotrino/fusion/spec/app"
)

type Repository[T any] struct {
	Path string // the resource path like /api/v1/books
}

func (Repository[T]) IsRepository() bool {
	return true
}

func (r Repository[T]) New(ctx context.Context) app.RepositoryImplStencil {
	return rest.REST[T](ctx, r.Path).ToStencil()
}

type Resource[T any] struct {
	Path    string // the resource path like /api/v1/book/42
	Default T      // a default value to use for populating an empty entity, e.g. for creation.
}

func (Resource[T]) IsResource() bool {
	return true
}

func (r Resource[T]) GetDefault() any {
	return r.Default
}

func (r Resource[T]) New(ctx context.Context) app.ResourceImplStencil {
	return rest.NewResource[T](ctx, r.Path).ToStencil()
}
