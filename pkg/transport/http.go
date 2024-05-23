package transport

import (
	"context"
	"strings"
)

/*
Se encarga de manejar las peticiones HTTP, decodificarlas y codificarlas.
*/

import (
	"net/http"
)

type Transport interface {
	Server(
		endpoint EndPoint,
		decode func(ctx context.Context, r *http.Request) (interface{}, error),
		encode func(ctx context.Context, w http.ResponseWriter, response interface{}) error,
		encodeError func(ctx context.Context, err error, w http.ResponseWriter),
	)
}

type EndPoint func(ctx context.Context, request interface{}) (interface{}, error)

type transport struct {
	w   http.ResponseWriter
	r   *http.Request
	ctx context.Context
}

func New(w http.ResponseWriter, r *http.Request, ctx context.Context) Transport {
	return &transport{
		w:   w,
		r:   r,
		ctx: ctx,
	}
}

func (t *transport) Server(
	endpoint EndPoint,
	decode func(ctx context.Context, r *http.Request) (interface{}, error),
	encode func(ctx context.Context, w http.ResponseWriter, response interface{}) error,
	encodeError func(ctx context.Context, err error, w http.ResponseWriter),
) {
	request, err := decode(t.ctx, t.r)
	if err != nil {
		encodeError(t.ctx, err, t.w)
		return
	}

	response, err := endpoint(t.ctx, request)
	if err != nil {
		encodeError(t.ctx, err, t.w)
		return
	}

	e := encode(t.ctx, t.w, response)
	if e != nil {
		encodeError(t.ctx, err, t.w)
		return
	}
}

func Clean(url string) ([]string, int) {
	if url[0] != '/' {
		url = "/" + url
	}

	if url[len(url)-1] != '/' {
		url = url + "/"
	}

	parts := strings.Split(url, "/")

	return parts, len(parts)
}
