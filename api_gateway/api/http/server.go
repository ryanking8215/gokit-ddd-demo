package http

import (
	"net/http"

	"gokit-ddd-demo/api_gateway/api"
	"gokit-ddd-demo/api_gateway/svc"
	"gokit-ddd-demo/lib/kitx"

	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/gorilla/mux"
)

func NewHTTPHandler(s svc.Service, opts ...kitx.Option) http.Handler {
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

	// options := []httptransport.ServerOption{
	// 	// httptransport.ServerErrorEncoder(errorEncoder),
	// 	//httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	// }

	r := mux.NewRouter()

	makeUserHTTPHandler(s, r)

	return r
}

func makeUserHTTPHandler(s svc.Service, r *mux.Router, opts ...kitx.Option) {
	{
		// ep := kitx.ServerEndpoint(func() (endpoint.Endpoint, string) {
		// 	ep := api.MakeFindUsersEndpoint(s)
		// 	return ep, "user_svc.Find"
		// }, opts...)

		r.Handle("/api/users", httptransport.NewServer(
			api.MakeFindUsersEndpoint(s),
			decodeFindUserRequest,
			encodeResponse,
			//options...,
		//append(options, httptransport.ServerBefore(opentracing.HTTPToContext(otTracer, "Sum", logger)))...,
		)).Methods("GET")
	}

	// TODO
	// r.Handle("/users/{id:[0-9]+}", httptransport.NewServer(
	// 	api.MakeGetUserEndpoint(s),
	// 	decodeGetRequest,
	// 	encodeResponse,
	// 	options...,
	// //append(options, httptransport.ServerBefore(opentracing.HTTPToContext(otTracer, "Concat", logger)))...,
	// )).Methods("GET")
}
