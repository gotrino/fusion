package rest

// A Repository which represents CRUD (create read update delete) operations on an Entity based resource set.
// If any entity has a certain kind of ID, the repository implementation must unmarshal it from a string
// to support Load and Delete.
type Repository[T any] interface {
	List() ([]T, error)
	Load(id string) (T, error)
	// Delete removes the entity. It is no error, if an already deleted entry is removed again.
	Delete(id string) error
	// Save updates or creates the Entity.
	Save(t T) error
}

// ResourceRepository represents an aggregate which may or may not have an id.
type ResourceRepository[T any] interface {
	Load() (T, error)
	Save(t T) error
	Delete() error
}
