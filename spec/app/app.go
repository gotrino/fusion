package app

import (
	"context"
)

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

func Navigate(ctx context.Context, params ActivityComposer) {

}

// An Activity declares a bunch of Fragments.
type Activity struct {
	Title     string
	Visible   bool
	Launcher  Launcher
	Fragments []Fragment
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

// Repository is a marker interface for a repository specification which represents a collection of resources.
type Repository interface {
	IsRepository() bool
}

// Resource is a marker interface to represent a single collection.
type Resource interface {
	IsResource() bool
}
