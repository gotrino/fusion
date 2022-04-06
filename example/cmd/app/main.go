package main

import (
	"context"
	"github.com/gotrino/fusion/example/turbines"
	"github.com/gotrino/fusion/runtime"
	"github.com/gotrino/fusion/spec/app"
	_ "github.com/gotrino/fusion/wasmjs"
)

type MyApp struct {
}

func (a MyApp) Compose(ctx context.Context) app.Application {
	return app.Application{
		Activities: []app.ActivityComposer{
			turbines.Overview{},
			turbines.Details{},
		},
		Authentication: app.Bearer{},
	}
}

func main() {
	runtime.MustStart("wasm/js", MyApp{})
}
