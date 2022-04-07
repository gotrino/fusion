package table

import (
	"context"
	"github.com/gotrino/fusion/spec/app"
)

type DataTableStencil struct {
	Repository app.Repository
	Deletable  bool
	Columns    []string
	OnRender   func(ctx context.Context, item any, col int) string
	OnClick    func(ctx context.Context, item any)
}

type DataTable[T any] struct {
	Repository app.Repository
	Deletable  bool
	Columns    []string
	OnRender   func(ctx context.Context, item T, col int) string
	OnClick    func(ctx context.Context, item T)
}

func (DataTable[T]) IsFragment() bool {
	return true
}

func (t DataTable[T]) ToStencil() any {
	return DataTableStencil{
		Repository: t.Repository,
		Deletable:  t.Deletable,
		Columns:    t.Columns,
		OnRender: func(ctx context.Context, item any, col int) string {
			if t.OnRender != nil {
				return t.OnRender(ctx, item.(T), col)
			}

			return ""
		},
		OnClick: func(ctx context.Context, item any) {
			if t.OnClick != nil {
				t.OnClick(ctx, item.(T))
			}
		},
	}
}
