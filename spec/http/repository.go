package http

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gotrino/fusion/spec/app"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
)

const (
	DecoderError = -1
	EncoderError = -2
)

type Request = http.Request
type Response = http.Response

const (
	StatusOK        = http.StatusOK
	StatusCreated   = http.StatusCreated
	StatusAccepted  = http.StatusAccepted
	StatusNoContent = http.StatusNoContent
)

type Repository[T any] struct {
	OnSave   func(t T) error
	OnLoad   func(id string) (T, error)
	OnList   func() ([]T, error)
	OnDelete func(id string) error
}

func (r Repository[T]) List() ([]any, error) {
	if r.OnList == nil {
		return nil, fmt.Errorf("OnSave is not implemented")
	}

	res, err := r.OnList()
	if err != nil {
		return nil, err
	}

	boxed := make([]any, 0, len(res))
	for _, t := range res {
		boxed = append(boxed, t)
	}

	return boxed, nil
}

func (r Repository[T]) Delete(id string) error {
	if r.OnDelete == nil {
		return fmt.Errorf("OnDelete is not implemented")
	}

	return r.OnDelete(id)
}

func (r Repository[T]) Save(entity any) error {
	if r.OnSave == nil {
		return fmt.Errorf("OnSave is not implemented")
	}

	return r.OnSave(entity.(T))
}

func (r Repository[T]) Load(id string) (any, error) {
	if r.OnLoad == nil {
		return nil, fmt.Errorf("OnLoad is not implemented")
	}

	return r.OnLoad(id)
}

func (Repository[T]) GetDefault() any {
	var t T
	return t
}
func (r Repository[T]) New(ctx context.Context) app.RepositoryImplStencil {
	return r
}

func (Repository[T]) IsRepository() bool {
	return true
}

func NewRequest(ctx context.Context, method string, url *url.URL, body io.Reader) *Request {
	req, err := http.NewRequestWithContext(ctx, method, url.String(), body)
	if err != nil {
		panic(fmt.Errorf("cannot happen: %w", err))
	}

	return req
}

func URL(ctx context.Context, paths ...string) *url.URL {
	a := app.FromContext[app.Application](ctx)
	p := path.Join(paths...)
	if strings.HasPrefix(p, "/") {
		p = p[1:]
	}

	base, err := url.Parse(fmt.Sprintf("%s://%s:%d/%s", a.Connection.Scheme, a.Connection.Host, a.Connection.Port, p))
	if err != nil {
		panic(fmt.Errorf("invalid url: %w", err))
	}

	return base
}

func Client(ctx context.Context) *http.Client {
	return http.DefaultClient
}

type Params struct {
	ContentType string
	Body        []byte
}

func Do(ctx context.Context, method string, url *url.URL, params Params, acceptableStatus ...int) ([]byte, error) {
	var body io.Reader
	if params.Body != nil {
		body = bytes.NewReader(params.Body)
	}
	req := NewRequest(ctx, method, url, body)
	if params.ContentType != "" {
		req.Header.Set("Content-Type", params.ContentType)
	}

	req = Authorizer(ctx)(req)

	res, err := Client(ctx).Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	log.Printf("%+v", res)

	statusFound := false
	for _, status := range acceptableStatus {
		if res.StatusCode == status {
			statusFound = true
			break
		}
	}

	if !statusFound {
		return nil, HttpError{Status: res.StatusCode}
	}

	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func Authorizer(ctx context.Context) func(request *http.Request) *http.Request {
	return func(request *http.Request) *http.Request {
		myApp := app.FromContext[app.Application](ctx)
		switch t := myApp.Authentication.(type) {
		case app.HardcodedBearer:
			request.Header.Set("Authorization", "Bearer "+t.Token)
		default:
			panic(fmt.Errorf("unsupported authorization: %T", t))
		}

		return request
	}
}

type HttpError struct {
	Status int
	Cause  error
}

func (e HttpError) Error() string {
	return fmt.Sprintf("http-error: %d", e.Status)
}

func (e HttpError) Forbidden() bool {
	return e.Status == http.StatusForbidden
}

func (e HttpError) Unauthenticated() bool {
	return e.Status == http.StatusUnauthorized
}

func (e HttpError) InternalServerError() bool {
	return e.Status >= http.StatusInternalServerError && e.Status <= http.StatusVariantAlsoNegotiates
}

func (e HttpError) ProtocolError() bool {
	return e.Status == DecoderError || e.Status == EncoderError
}

func (e HttpError) NotFound() bool {
	return e.Status == http.StatusNotFound
}

func (e HttpError) Unwrap() error {
	return e.Cause
}
