package turbines

import (
	"context"
	"github.com/gotrino/fusion/spec/app"
	"github.com/gotrino/fusion/spec/rest"
	"github.com/gotrino/fusion/spec/table"
)

type Overview struct {
	Route string `route:"/devices"` // gets injected, due to app.Route, matches the route-tag
	Limit int    `query:"limit"`    // parsed from the routing system, due to query-tag
}

func (a Overview) Compose(ctx context.Context) app.Activity {
	return app.Activity{
		Title: "Windturbinen",
		Launcher: app.Icon{
			Name:  "fa-settings",
			Title: "Turbinen zeigen",
			Hint:  "Zeigt alle Turbinen",
		},
		Visible: true,
		Fragments: []app.Fragment{
			table.DataTable[Turbine]{
				Repository: rest.Repository[Turbine]{
					Path: "/api/v1/turbines",
				},
				Deletable: true,
				Columns:   []string{"Name", "Leistung"},
				OnRender: func(ctx context.Context, t Turbine, col int) string {
					switch col {
					case 0:
						return t.ID
					case 1:
						return t.Name
					case 2:
						return "12kw/h"
					default:
						return ""
					}
				},
				OnClick: func(ctx context.Context, t Turbine) {
					app.Navigate(ctx, Details{ID: t.ID})
				},
			},
		},
	}
}
