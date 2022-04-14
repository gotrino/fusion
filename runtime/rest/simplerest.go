package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gotrino/fusion/spec/app"
	http2 "github.com/gotrino/fusion/spec/http"
	"log"
	"path"
	"strings"

	"io"
	"net/http"
	"net/url"
	"reflect"
)

func REST[T any](ctx context.Context, resource string) RESTRepo[T] {
	a := app.FromContext[app.Application](ctx)
	if strings.HasPrefix(resource, "/") {
		resource = resource[1:]
	}
	base, err := url.Parse(fmt.Sprintf("%s://%s:%d/%s", a.Connection.Scheme, a.Connection.Host, a.Connection.Port, resource))
	if err != nil {
		panic(err)
	}
	log.Println("!! rest repo using", base.String())

	return RESTRepo[T]{Base: base, WithRequest: http2.Authorizer(ctx)}
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

func (r RESTRepo[T]) ToStencil() app.RepositoryImplStencil {
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
		return nil, http2.HttpError{Status: resp.StatusCode}
	}

	var res []T
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&res); err != nil {
		return nil, http2.HttpError{Status: http2.DecoderError, Cause: err}
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
		return res, http2.HttpError{Status: resp.StatusCode}
	}

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&res); err != nil {
		return res, http2.HttpError{Status: http2.DecoderError, Cause: err}
	}

	return res, nil
}

// Delete performs a delete on the root resource attached with the id, like DELETE /api/movies/{id}.
func (r RESTRepo[T]) Delete(id string) error {
	req := r.req("DELETE", id, nil)
	log.Println(">>", req.URL)
	resp, err := r.client().Do(req)
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
		return http2.HttpError{Status: resp.StatusCode}
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
		return http2.HttpError{Status: http2.EncoderError, Cause: err}
	}

	req := r.req("PUT", id, bytes.NewReader(buf))
	req.Header.Set("Content-Type", "application/json")
	resp, err := r.client().Do(req)
	if err != nil {
		return http2.HttpError{Cause: err}
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusAccepted:
		fallthrough
	case http.StatusNoContent:
		fallthrough
	case http.StatusCreated:
		fallthrough
	case http.StatusOK:
		return nil
	default:
		return http2.HttpError{Status: resp.StatusCode}
	}
}

func (r RESTRepo[T]) client() *http.Client {
	if r.Client == nil {
		return http.DefaultClient
	}

	return r.Client
}

func (r RESTRepo[T]) req(method string, p string, body io.Reader) *http.Request {
	if r.Base == nil {
		u, err := url.Parse("http://localhost:8080")
		if err != nil {
			panic(err)
		}

		r.Base = u
	}

	u, err := r.Base.Parse(path.Join(r.Base.String(), p))
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
