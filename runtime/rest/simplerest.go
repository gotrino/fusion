package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gotrino/fusion/spec/app"

	"io"
	"net/http"
	"net/url"
	"reflect"
)

type RepositoryImplStencil interface {
	List() ([]any, error)        // any is of type []interface boxing each T
	Load(id string) (any, error) // any is of type T
	Delete(id string) error
	Save(t any) error // any is of type T
}

func REST[T any](ctx context.Context, resource string) RESTRepo[T] {
	base, err := url.Parse("http://localhost:8081" + resource)
	if err != nil {
		panic(err)
	}

	myApp := app.FromContext[app.Application](ctx)

	return RESTRepo[T]{Base: base, WithRequest: func(request *http.Request) *http.Request {
		switch t := myApp.Authentication.(type) {
		case app.HardcodedBearer:
			request.Header.Set("Authorization", "Bearer "+t.Token)
		default:
			panic(fmt.Errorf("unsupported authorization: %T", t))
		}

		return request
	}}
}

// RESTRepo is a simple more or less idiomatic REST based CRUD repository adapter. It makes really strong assumptions
// about the verbs.
type RESTRepo[T any] struct {
	Context     context.Context
	Base        *url.URL
	WithRequest func(*http.Request) *http.Request
	Client      *http.Client
	// the actual resource like /api/movie
	Resource string
}

func (r RESTRepo[T]) ToStencil() RepositoryImplStencil {
	return stencilAdapter[T]{r}
}

// List performs a get on the root resource, like GET /api/movies and expects a json array.
func (r RESTRepo[T]) List() ([]T, error) {
	resp, err := r.client().Do(r.req("GET", "", nil))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, httpError{status: resp.StatusCode}
	}

	var res []T
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&res); err != nil {
		return nil, httpError{status: decoderError, cause: err}
	}

	return res, nil
}

// Load performs a get on the root resource attached with the id, like GET /api/movies/{id}.
func (r RESTRepo[T]) Load(id string) (T, error) {
	var res T
	resp, err := r.client().Do(r.req("GET", id, nil))
	if err != nil {
		return res, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return res, httpError{status: resp.StatusCode}
	}

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&res); err != nil {
		return res, httpError{status: decoderError, cause: err}
	}

	return res, nil
}

// Delete performs a delete on the root resource attached with the id, like DELETE /api/movies/{id}.
func (r RESTRepo[T]) Delete(id string) error {
	resp, err := r.client().Do(r.req("DELETE", id, nil))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusAccepted:
		fallthrough
	case http.StatusNoContent:
		fallthrough
	case http.StatusOK:
	default:
		return httpError{status: resp.StatusCode}
	}

	return nil
}

// Save performs a put on the root resource attached with the id, like PUT /api/movies/{id}.
func (r RESTRepo[T]) Save(t T) error {
	id, err := GetID(t)
	if err != nil {
		panic(err)
	}

	buf, err := json.Marshal(t)
	if err != nil {
		return httpError{status: encoderError, cause: err}
	}

	resp, err := r.client().Do(r.req("PUT", id, bytes.NewReader(buf)))
	if err != nil {
		return nil
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusAccepted:
		fallthrough
	case http.StatusNoContent:
		fallthrough
	case http.StatusOK:
		return nil
	default:
		return httpError{status: resp.StatusCode}
	}
}

func (r RESTRepo[T]) client() *http.Client {
	if r.Client == nil {
		return http.DefaultClient
	}

	return r.Client
}

func (r RESTRepo[T]) req(method string, path string, body io.Reader) *http.Request {
	if r.Base == nil {
		u, err := url.Parse("http://localhost:8080")
		if err != nil {
			panic(err)
		}

		r.Base = u
	}

	u, err := r.Base.Parse(path)
	if err != nil {
		panic(err)
	}

	ctx := r.Context
	if ctx == nil {
		ctx = context.Background()
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		panic(err)
	}

	if r.WithRequest != nil {
		req = r.WithRequest(req)
	}

	return req
}

// GetID returns the entity either by calling ID()string method or grabbing the ID string field.
func GetID(a any) (string, error) {
	if ider, ok := a.(interface{ ID() string }); ok {
		return ider.ID(), nil
	}

	t := reflect.TypeOf(a)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Name == "ID" && f.Type.Name() == "string" {
			return reflect.ValueOf(a).Field(i).String(), nil
		}
	}

	return "", fmt.Errorf("type " + t.String() + " must either provider 'ID() string' method oder 'ID string' field")
}

type stencilAdapter[T any] struct {
	impl RESTRepo[T]
}

func (s stencilAdapter[T]) List() ([]any, error) {
	res, err := s.impl.List()
	if err != nil {
		return nil, err
	}

	boxed := make([]any, 0, len(res))
	for _, t := range res {
		boxed = append(boxed, t)
	}

	return boxed, nil
}

func (s stencilAdapter[T]) Load(id string) (any, error) {
	return s.impl.Load(id)
}

func (s stencilAdapter[T]) Delete(id string) error {
	return s.impl.Delete(id)
}

func (s stencilAdapter[T]) Save(t any) error {
	return s.impl.Save(t.(T))
}
