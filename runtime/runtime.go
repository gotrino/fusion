package runtime

import (
	"context"
	"fmt"
	"github.com/gotrino/fusion/spec/app"
	"sync"
)

var runtimes = map[string]Factory{}
var lock sync.Mutex

type State struct {
	Context     context.Context
	Application app.Application
	Activities  []app.Activity
	Active      int
}

type Runtime interface {
	Start(spec app.ApplicationComposer) error
}

type Factory func() (Runtime, error)

func Register(name string, factory Factory) {
	lock.Lock()
	defer lock.Unlock()

	runtimes[name] = factory
}

func Open(name string) (Runtime, error) {
	lock.Lock()
	defer lock.Unlock()

	fac, ok := runtimes[name]
	if !ok {
		return nil, fmt.Errorf("runtime '%s' is not available", name)
	}

	rt, err := fac()
	if err != nil {
		return nil, fmt.Errorf("cannot create an instance of runtime '%s': %w", name, err)
	}

	return rt, nil
}

func MustStart(name string, spec app.ApplicationComposer) {
	rt, err := Open(name)
	if err != nil {
		panic(err)
	}

	if err := rt.Start(spec); err != nil {
		panic(err)
	}
}
