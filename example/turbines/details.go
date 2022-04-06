package turbines

import (
	"context"
	"github.com/gotrino/fusion/spec/app"
	"github.com/gotrino/fusion/spec/form"
	"github.com/gotrino/fusion/spec/rest"
	"strings"
)

type Details struct {
	Route app.Route `route:"/devices/:id/details"`
	ID    string    `path:"id"`
}

func (a Details) Compose(ctx context.Context) app.Activity {
	return app.Activity{
		Fragments: []app.Fragment{
			form.Form[Turbine]{
				Resource: rest.Resource[Turbine]{
					Path: "/api/v1/turbines/" + a.ID,
				},
				Fields: []form.Field{
					form.Header{
						Text: "Turbine details",
						Hint: "a funny turbine thing...",
					},
					form.Text[Turbine]{
						Text: "Name",
						Hint: "the name you shout if it breaks",
						ToModel: func(src string, dst *Turbine) error {
							if strings.HasPrefix("wdy-", src) {
								return app.ValidationError{Message: "should start with wdy- prefix"}
							}

							dst.Name = src
							return nil
						},
						FromModel: func(src Turbine) string {
							return src.Name
						},
					},
					form.Integer[Turbine]{
						Text:     "Degree",
						Disabled: true,
						FromModel: func(src Turbine) int64 {
							return 100
						},
					},
					form.Select[Turbine]{
						Text:        "Select something",
						Hint:        "a cool selector",
						MultiSelect: true,
						ToModel: func(src []form.Item, dst *Turbine) error {
							return nil
						},
						FromModel: func(src Turbine) []form.Item {
							return nil
						},
					},
				},
			},
		},
	}
}
