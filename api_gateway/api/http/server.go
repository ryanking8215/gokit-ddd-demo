package http

import (
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/tracing/zipkin"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"gokit-ddd-demo/api_gateway/api"
	"gokit-ddd-demo/api_gateway/svc"
	"gokit-ddd-demo/lib/kitx"
)

func NewHTTPHandler(s svc.Service, opts *kitx.ServerOptions) http.Handler {
	// Zipkin HTTP Server Trace can either be instantiated per endpoint with a
	// provided operation name or a global tracing service can be instantiated
	// without an operation name and fed to each Go kit endpoint as ServerOption.
	// In the latter case, the operation name will be the endpoint's http method.
	// We demonstrate a global tracing service here.

	tracer := opts.ZipkinTracer()
	zipkinServer := zipkin.HTTPServerTrace(tracer)

	logger := opts.Logger()

	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(errorEncoder),
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		zipkinServer,
	}

	r := mux.NewRouter()

	makeUserHTTPHandler(s, r, options, opts)

	return r
}

func makeUserHTTPHandler(s svc.Service, r *mux.Router, httpOpts []httptransport.ServerOption, opts *kitx.ServerOptions) {
	{
		ep := kitx.ServerEndpoint(func() (endpoint.Endpoint, string) {
			ep := api.MakeFindUsersEndpoint(s)
			return ep, "user_svc.Find"
		}, opts)
		r.Handle("/api/users", httptransport.NewServer(
			ep,
			decodeFindUserRequest,
			encodeResponse,
			httpOpts...,
		)).Methods("GET")
	}

	{
		ep := kitx.ServerEndpoint(func() (endpoint.Endpoint, string) {
			ep := api.MakeGetUserEndpoint(s)
			return ep, "user_svc.Get"
		}, opts)
		r.Handle("/api/users/{id:[0-9]+}", httptransport.NewServer(
			ep,
			decodeGetUserRequest,
			encodeResponse,
			httpOpts...,
		//append(options, httptransport.ServerBefore(opentracing.HTTPToContext(otTracer, "Concat", logger)))...,
		)).Methods("GET")
	}
}
