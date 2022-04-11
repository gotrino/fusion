package table

import (
	"context"
	"github.com/gotrino/fusion/spec/app"
	"github.com/gotrino/fusion/spec/svg"
)

type DataTableStencil struct {
	Repository app.Repository
	Deletable  bool
	Columns    []Column
	OnRender   func(ctx context.Context, item any, col int) Cell
	OnClick    func(ctx context.Context, item any)
}

type Cell struct {
	Values     []string
	RenderHint string
}

func NewText(v string) Cell {
	return Cell{
		Values:     []string{v},
		RenderHint: "text-1",
	}
}

func NewSVG(svg svg.SVG, text ...string) Cell {
	for len(text) < 2 {
		text = append(text, "")
	}

	return Cell{
		Values:     []string{text[0], text[1], string(svg)},
		RenderHint: "svg-text-2",
	}
}

type Column struct {
	Name   string
	Weight int
}

type DataTable[T any] struct {
	Repository app.Repository
	Deletable  bool
	Columns    []Column
	OnRender   func(ctx context.Context, item T, col int) Cell
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
		OnRender: func(ctx context.Context, item any, col int) Cell {
			if t.OnRender != nil {
				return t.OnRender(ctx, item.(T), col)
			}

			return Cell{}
		},
		OnClick: func(ctx context.Context, item any) {
			if t.OnClick != nil {
				t.OnClick(ctx, item.(T))
			}
		},
	}
}
