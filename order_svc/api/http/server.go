package http

import (
	"net/http"
	//"time"

	//"github.com/go-kit/kit/tracing/opentracing"
	//"github.com/go-kit/kit/tracing/zipkin"
	"github.com/go-kit/kit/log"
	//"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	apiendpoint "gokit-ddd-demo/order_svc/api/endpoint"
	"gokit-ddd-demo/order_svc/domain/order"
)

func NewHTTPHandler(svc order.Service, logger log.Logger) http.Handler {
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
		ep := apiendpoint.MakeFindEndpoint(svc)
		ep = apiendpoint.MiscMiddleware("ordersvc.Find", logger, nil)(ep)
		m.Handle("/orders", httptransport.NewServer(
			ep,
			decodeFindRequest,
			encodeResponse,
			options...,
		//append(options, httptransport.ServerBefore(opentracing.HTTPToContext(otTracer, "Sum", logger)))...,
		)).Methods("GET")
	}
	{
		ep := apiendpoint.MakeGetEndpoint(svc)
		ep = apiendpoint.MiscMiddleware("ordersvc.Get", logger, nil)(ep)
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
