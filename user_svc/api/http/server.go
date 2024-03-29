package http

import (
	"net/http"

	"gokit-ddd-demo/lib/kitx"
	"gokit-ddd-demo/user_svc/api"
	"gokit-ddd-demo/user_svc/svc/user"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/gorilla/mux"
)

func NewHTTPHandler(svc user.Service, opts *kitx.ServerOptions) http.Handler {
	// Zipkin HTTP Server Trace can either be instantiated per endpoint with a
	// provided operation name or a global tracing service can be instantiated
	// without an operation name and fed to each Go kit endpoint as ServerOption.
	// In the latter case, the operation name will be the endpoint's http method.
	// We demonstrate a global tracing service here.
	//zipkinServer := zipkin.HTTPServerTrace(zipkinTracer)
	//
	//options := []httptransport.ServerOption{
	//	httptransport.ServerErrorEncoder(errorEncoder),
	//	httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	//	zipkinServer,
	//}
	logger := opts.Logger()

	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(errorEncoder),
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	m := mux.NewRouter()
	{

		ep := kitx.ServerEndpoint(func() (endpoint.Endpoint, string) {
			ep := api.MakeFindEndpoint(svc)
			return ep, "user_svc.Find"
		}, opts)

		m.Handle("/users", httptransport.NewServer(
			ep,
			decodeFindRequest,
			encodeResponse,
			options...,
		// append(options, httptransport.ServerBefore(opentracing.HTTPToContext(otTracer, "Sum", logger)))...,
		)).Methods("GET")
	}
	{
		ep := kitx.ServerEndpoint(func() (endpoint.Endpoint, string) {
			ep := api.MakeGetEndpoint(svc)
			return ep, "user_svc.Get"
		}, opts)

		m.Handle("/users/{id:[0-9]+}", httptransport.NewServer(
			ep,
			decodeGetRequest,
			encodeResponse,
			options...,
		//append(options, httptransport.ServerBefore(opentracing.HTTPToContext(otTracer, "Concat", logger)))...,
		)).Methods("GET")
	}

	return m
}
