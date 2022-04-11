package rest

import (
	"context"
	"github.com/gotrino/fusion/runtime/rest"
)

type RepositoryImplStencil interface {
	List() ([]any, error)        // any is of type []T
	Load(id string) (any, error) // any is of type T
	Delete(id string) error
	Save(t any) error // any is of type T
}

type Repository[T any] struct {
	Path string // the resource path like /api/v1/books
}

func (Repository[T]) IsRepository() bool {
	return true
}

func (r Repository[T]) New(ctx context.Context) RepositoryImplStencil {
	return rest.REST[T](ctx, r.Path).ToStencil()
}

type Resource[T any] struct {
	Path    string // the resource path like /api/v1/book/42
	Default T      // a default value to use for populating an empty entity
}

func (Resource[T]) IsResource() bool {
	return true
}
