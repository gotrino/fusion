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
	"strings"
)

// TODO this has a lot of duplication which should be unified.
type RestResourceImpl[T any] struct {
	Context     context.Context
	Base        *url.URL
	WithRequest func(*http.Request) *http.Request
	Client      *http.Client
}

func NewResource[T any](ctx context.Context, resource string) RestResourceImpl[T] {
	a := app.FromContext[app.Application](ctx)
	if strings.HasPrefix(resource, "/") {
		resource = resource[1:]
	}
	base, err := url.Parse(fmt.Sprintf("%s://%s:%d/%s", a.Connection.Scheme, a.Connection.Host, a.Connection.Port, resource))
	if err != nil {
		panic(err)
	}

	return RestResourceImpl[T]{Base: base, Context: ctx, WithRequest: func(request *http.Request) *http.Request {
		switch t := a.Authentication.(type) {
		case app.HardcodedBearer:
			request.Header.Set("Authorization", "Bearer "+t.Token)
		default:
			panic(fmt.Errorf("unsupported authorization: %T", t))
		}

		return request
	}}
}

func (r RestResourceImpl[T]) Load() (T, error) {
	var res T
	resp, err := r.client().Do(r.req("GET", "", nil))
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

func (r RestResourceImpl[T]) Save(t T) error {
	id, err := GetID(t)
	if err != nil {
		panic(err)
	}

	buf, err := json.Marshal(t)
	if err != nil {
		return httpError{status: encoderError, cause: err}
	}

	req := r.req("PUT", id, bytes.NewReader(buf))
	req.Header.Set("Content-Type", "application/json")
	resp, err := r.client().Do(req)

	if err != nil {
		return httpError{cause: err}
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
		return httpError{status: resp.StatusCode}
	}
}

func (r RestResourceImpl[T]) Delete() error {
	resp, err := r.client().Do(r.req("DELETE", "", nil))
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

func (r RestResourceImpl[T]) client() *http.Client {
	if r.Client == nil {
		return http.DefaultClient
	}

	return r.Client
}

func (r RestResourceImpl[T]) req(method string, path string, body io.Reader) *http.Request {
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

func (r RestResourceImpl[T]) ToStencil() app.ResourceImplStencil {
	return stencilResAdapter[T]{impl: r}
}

type stencilResAdapter[T any] struct {
	impl RestResourceImpl[T]
}

func (s stencilResAdapter[T]) Load() (any, error) {
	return s.impl.Load()
}

func (s stencilResAdapter[T]) Delete() error {
	return s.impl.Delete()
}

func (s stencilResAdapter[T]) Save(t any) error {
	return s.impl.Save(t.(T))
}
