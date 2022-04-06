package rest

type Repository[T any] struct {
	Path    string // the resource path like /api/v1/books
	Default T      // a default value to use for populating an empty entity
}

func (Repository[T]) IsRepository() bool {
	return true
}

type Resource[T any] struct {
	Path    string // the resource path like /api/v1/book/42
	Default T      // a default value to use for populating an empty entity
}

func (Resource[T]) IsResource() bool {
	return true
}
