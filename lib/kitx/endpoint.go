package kitx

import (
	// "context"
	// "fmt"
	// "time"

	"gokit-ddd-demo/lib"
	"io"
	"sync"
	"time"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/kit/tracing/zipkin"
	"github.com/sony/gobreaker"
	"google.golang.org/grpc"
)

func ServerEndpoint(makeEndpoint func() (endpoint.Endpoint, string), options *ServerOptions) endpoint.Endpoint {
	ep, name := makeEndpoint()

	if options.rateLimitOption.limiter != nil {
		ep = ratelimit.NewErroringLimiter(options.rateLimitOption.limiter)(ep)
	}
	if options.circuitBreakerOption.enable {
		ep = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(ep)
	}
	if options.openTracingOption.otTracer != nil {
		ep = opentracing.TraceServer(options.openTracingOption.otTracer, name)(ep)
	}
	if options.zipkinTracerOption.tracer != nil {
		ep = zipkin.TraceEndpoint(options.zipkinTracerOption.tracer, name)(ep)
	}
	if options.Logger() != nil {
		ep = LoggingMiddleware(name, log.With(options.Logger(), "method", name))(ep)
	}
	if options.metricsOption.histogram != nil {
		ep = InstrumentingMiddleware(options.metricsOption.histogram.With("method", name))(ep)
	}

	return ep
}

var GRPCConnections sync.Map

func newGRPCClientFactory(makeEndpoint func(conn *grpc.ClientConn) (endpoint.Endpoint, string), opts *ClientOptions) sd.Factory {
	return func(instance string) (i endpoint.Endpoint, closer io.Closer, e error) {
		ref := NewGRPConnRef()
		actual, _ := GRPCConnections.LoadOrStore(instance, ref)
		ref = actual.(*GRPCConnRef)

		conn, err := ref.Conn(instance)
		if err != nil {
			return nil, nil, err
		}
		ep, name := makeEndpoint(conn)

		if opts.openTracingOption.otTracer != nil {
			ep = opentracing.TraceClient(opts.openTracingOption.otTracer, name)(ep)
		}

		// if opts.rateLimitOption.enable {
		// 	ep = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(opts.rateLimitOption.every), opts.rateLimitOption.tokenCnt))(ep)
		// }

		// if opts.circuitBreakerOption.enable {
		// 	ep = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
		// 		Name:    name,
		// 		Timeout: opts.circuitBreakerOption.timeout,
		// 	}))(ep)
		// }

		return ep, ref, nil
	}
}

func GRPCClientEndpoint(instancer sd.Instancer, makeEndpoint func(conn *grpc.ClientConn) (endpoint.Endpoint, string), opts *ClientOptions) endpoint.Endpoint {
	factory := newGRPCClientFactory(makeEndpoint, opts)

	logger := opts.Logger()
	endpointer := sd.NewEndpointer(instancer, factory, logger)
	balancer := lb.NewRoundRobin(endpointer)

	retryMax := 1
	timeout := 3 * time.Second
	if opts.loadBalanceOption.retryMax > 0 {
		retryMax = opts.loadBalanceOption.retryMax
		timeout = opts.loadBalanceOption.timeout
	}

	return lb.RetryWithCallback(timeout, balancer, func(n int, received error) (bool, error) {
		if _, ok := received.(lib.Error); ok {
			return false, received
		}
		return n < retryMax, nil
	})
}
