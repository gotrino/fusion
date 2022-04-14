package form

// CodeEditor represents a view which allows editing markup or coding text using
// line numbers and syntax highlighting.
type CodeEditor[T any] struct {
	Lang      string                             // Lang is a hint for the syntax highlighter, e.g. golang or json
	ReadOnly  bool                               // ReadOnly determines whether the view allows editing
	ToModel   func(src string, dst T) (T, error) // ToModel is a generic func callback to map from the view-model (a string) into the domain model.
	FromModel func(src T) string                 // FromModel is a generic func callback to map from the domain model into the view-model.
}

// GetToModel returns a stenciled version of ToModel.
func (e CodeEditor[T]) GetToModel() func(src string, dst any) (any, error) {
	if e.ToModel == nil {
		return nil
	}

	return func(src string, dst any) (any, error) {
		return e.ToModel(src, dst.(T))
	}
}

// GetFromModel returns a stenciled version of FromModel.
func (e CodeEditor[T]) GetFromModel() func(src any) string {
	if e.FromModel == nil {
		return nil
	}

	return func(src any) string {
		return e.FromModel(src.(T))
	}
}

func (e CodeEditor[T]) GetLang() string {
	return e.Lang
}

func (e CodeEditor[T]) IsReadOnly() bool {
	return e.ReadOnly
}

func (CodeEditor[T]) IsField() bool {
	return true
}
