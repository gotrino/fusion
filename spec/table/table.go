package table

import (
	"context"
	"github.com/gotrino/fusion/spec/app"
)

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
