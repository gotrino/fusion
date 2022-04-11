package form

import (
	"github.com/gotrino/fusion/spec/app"
)

type Form struct {
	Title       string
	Description string

	Resource app.Resource
	Fields   []Field
}

func (f Form) IsFragment() bool {
	return true
}

type StencilText struct {
	Label        string
	Description  string
	Disabled     bool
	DefaultValue string
	Lines        int
	ToModel      func(src string, dst any) error
	FromModel    func(src any) string
}

// Field is a marker interface to identify a field mapping.
type Field interface {
	IsField() bool
}

// Header is a kind of content separator for logical sections, which is not really an interactive but informative field.
type Header struct {
	Text string
	Hint string
}

func (Header) IsField() bool {
	return true
}

// A Text is something like an edit text field.
type Text[T any] struct {
	Label        string
	Description  string
	Disabled     bool
	DefaultValue string
	Lines        int
	ToModel      func(src string, dst *T) error
	FromModel    func(src T) string
}

func (t Text[T]) ToStencil() StencilText {
	return StencilText{
		Label:        t.Label,
		Description:  t.Description,
		Disabled:     t.Disabled,
		DefaultValue: t.DefaultValue,
		Lines:        t.Lines,
		ToModel: func(src string, dst any) error {
			x := dst.(T) // TODO not correct
			return t.ToModel(src, &x)
		},
		FromModel: func(src any) string {
			return t.FromModel(src.(T))
		},
	}
}

func (Text[T]) IsField() bool {
	return true
}

// An Integer field just allows per-se only integer numbers.
type Integer[T any] struct {
	Text      string
	Hint      string
	Disabled  bool
	ToModel   func(src int64, dst *T) error
	FromModel func(src T) int64
}

func (Integer[T]) IsField() bool {
	return true
}

type Select[T any] struct {
	Text        string
	Hint        string
	Disabled    bool
	MultiSelect bool
	Format      string // Format is a render hint like combobox or radiobutton or checkbox. This may be ignored by the renderer.
	ToModel     func(src []Item, dst *T) error
	FromModel   func(src T) []Item
}

func (Select[T]) IsField() bool {
	return true
}

type Item struct {
	ID       string // not shown, may be empty - just for the developers convenience
	Text     string // text to show
	Hint     string // eventually shows the hint
	Selected bool
}
