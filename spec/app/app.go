package app

import (
	"context"
	"fmt"
	"log"
)

type Navigator struct {
	Delegate interface {
		Navigate(params ActivityComposer)
	}
}

func (n Navigator) Navigate(params ActivityComposer) {
	n.Delegate.Navigate(params)
}

// Fragment is a marker interface to identify composables of gotrino zero.
type Fragment interface {
	IsFragment() bool
}

// Launcher is a marker interface for entry points. An application must have at least a single activity with
// a launcher.
type Launcher interface {
	IsLauncher() bool
}

// A Route provides some marshall and unmarshal logic.
type Route string

// Navigate assembles a query link based on the given composer params, to ease things.
func Navigate(ctx context.Context, params ActivityComposer) {
	FromContext[Navigator](ctx).Navigate(params)
}

// An Activity declares a bunch of Fragments.
type Activity struct {
	Title     string
	Visible   bool
	Launcher  Launcher
	Fragments []Fragment
}

type Connection struct {
	Scheme string
	Host   string
	Port   int
}

// ActivityComposer creates and describes a concrete Activity instance.
type ActivityComposer interface {
	Compose(ctx context.Context) Activity
}

// Application declares an Application and its Activities or Pages - depending on the actual renderer.
type Application struct {
	Title          string
	Activities     []ActivityComposer
	Authentication Authentication
	Connection     Connection
}

// An ApplicationComposer creates and describes a concrete Application instance.
type ApplicationComposer interface {
	Compose(ctx context.Context) Application
}

type Repo struct {
}

type RepositoryComposer interface {
	Compose(ctx context.Context)
}

type RepositoryImplStencil interface {
	List() ([]any, error)        // any is of type []T
	Load(id string) (any, error) // any is of type T
	Delete(id string) error
	Save(t any) error // any is of type T
}

// Repository is a marker interface for a repository specification which represents a collection of resources.
type Repository interface {
	IsRepository() bool
	New(ctx context.Context) RepositoryImplStencil
}

type ResourceImplStencil interface {
	Load() (any, error) // any is of type T
	Delete() error
	Save(t any) error // any is of type T
}

// Resource is a marker interface to represent a single collection.
type Resource interface {
	New(ctx context.Context) ResourceImplStencil
	GetDefault() any
	IsResource() bool
}

type myCtxKey string

// FromContext cannot be used with interfaces because they boil down to any without type information.
func FromContext[T any](ctx context.Context) T {
	var t T
	k := fmt.Sprintf("%T", t)
	log.Printf("context: grabbing %s\n", k)

	a := ctx.Value(myCtxKey(k))
	return a.(T)
}

// WithContext cannot be used interfaces because they loose (any) type information.
func WithContext[T any](ctx context.Context, t T) context.Context {
	k := fmt.Sprintf("%T", t)
	log.Printf("context: setting %s\n", k)

	return context.WithValue(ctx, myCtxKey(k), t)
}
