package kitx

import (
	// "context"
	// "fmt"
	// "time"

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
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
)

func ServerEndpoint(makeEndpoint func() (endpoint.Endpoint, string), opts ...Option) endpoint.Endpoint {
	ep, name := makeEndpoint()

	options := options{}
	for _, o := range opts {
		o.apply(&options)
	}

	if options.rateLimitOption.enable {
		ep = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(options.rateLimitOption.every), options.rateLimitOption.tokenCnt))(ep)
	}
	if options.circuitBreakerOption.enable {
		ep = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(ep)
	}
	if options.openTracingOption.otTracer != nil {
		ep = opentracing.TraceServer(options.openTracingOption.otTracer, name)(ep)
	}

	// if zipkinTracer != nil {
	// 	ep = zipkin.TraceEndpoint(zipkinTracer, name)(ep)
	// }

	if options.loggerOption.logger != nil {
		ep = LoggingMiddleware(name, log.With(options.loggerOption.logger, "method", name))(ep)
	}

	// if options.metricsOption.enable {
	// 	ep = InstrumentingMiddleware(duration.With("method", name))(ep)
	// }

	return ep
}

var GRPCConnections sync.Map

func newGRPCClientFactory(makeEndpoint func(conn *grpc.ClientConn) (endpoint.Endpoint, string), opts *options) sd.Factory {
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

		if opts.rateLimitOption.enable {
			ep = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(opts.rateLimitOption.every), opts.rateLimitOption.tokenCnt))(ep)
		}

		if opts.circuitBreakerOption.enable {
			ep = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
				Name:    name,
				Timeout: opts.circuitBreakerOption.timeout,
			}))(ep)
		}

		return ep, ref, nil
	}
}

func GRPCClientEndpoint(instancer sd.Instancer, makeEndpoint func(conn *grpc.ClientConn) (endpoint.Endpoint, string), opts ...Option) endpoint.Endpoint {
	options := options{}
	for _, o := range opts {
		o.apply(&options)
	}

	factory := newGRPCClientFactory(makeEndpoint, &options)

	logger := options.loggerOption.logger
	endpointer := sd.NewEndpointer(instancer, factory, logger)
	balancer := lb.NewRoundRobin(endpointer)

	if options.loadBalanceOption.enable {
		return lb.Retry(options.loadBalanceOption.retryMax, options.loadBalanceOption.retryTimeout, balancer)
	}

	return lb.Retry(1, 3*time.Second, balancer)
}
