package http

import (
	"net/http"
	//"time"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"gokit-ddd-demo/lib/kitx"
	"gokit-ddd-demo/order_svc/api"
	"gokit-ddd-demo/order_svc/svc/order"
)

func NewHTTPHandler(svc order.Service, opts ...kitx.Option) http.Handler {
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

	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(errorEncoder),
		//httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	//m := http.NewServeMux()
	m := mux.NewRouter()
	{
		ep := kitx.ServerEndpoint(func() (endpoint.Endpoint, string) {
			ep := api.MakeFindEndpoint(svc)
			return ep, "order_svc.Find"
		}, opts...)
		m.Handle("/orders", httptransport.NewServer(
			ep,
			decodeFindRequest,
			encodeResponse,
			options...,
		//append(options, httptransport.ServerBefore(opentracing.HTTPToContext(otTracer, "Sum", logger)))...,
		)).Methods("GET")
	}
	{
		ep := kitx.ServerEndpoint(func() (endpoint.Endpoint, string) {
			ep := api.MakeGetEndpoint(svc)
			return ep, "order_svc.Get"
		}, opts...)
		m.Handle("/orders/{id:[0-9]+}", httptransport.NewServer(
			ep,
			decodeGetRequest,
			encodeResponse,
			options...,
		//append(options, httptransport.ServerBefore(opentracing.HTTPToContext(otTracer, "Concat", logger)))...,
		)).Methods("GET")
	}
	return m
}
