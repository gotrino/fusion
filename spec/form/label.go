package form

// A Label is a passive view containing arbitrary text derived by a concrete type.
// Line breaks should be respected.
type Label[T any] struct {
	Text      string             // Text is prefixed to the returned string from FromModel
	FromModel func(src T) string // FromModel is a generic func callback to map from the domain model into the view-model.
}

// GetFromModel returns a stenciled version of FromModel.
func (f Label[T]) GetFromModel() func(src any) string {
	if f.FromModel == nil {
		return nil
	}

	return func(src any) string {
		return f.FromModel(src.(T))
	}
}

func (f Label[T]) GetText() string {
	return f.Text
}

func (Label[T]) IsLabel() bool {
	return true
}

func (Label[T]) IsField() bool {
	return true
}
