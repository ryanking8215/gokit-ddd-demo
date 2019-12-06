package endpoint

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/ratelimit"
	//"github.com/go-kit/kit/tracing/opentracing"
	//"github.com/go-kit/kit/tracing/zipkin"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	//"gokit-ddd-demo/user_svc/domain/user"
)

// InstrumentingMiddleware returns an endpoint middleware that records
// the duration of each invocation to the passed histogram. The middleware adds
// a single field: "success", which is "true" if no error is returned, and
// "false" otherwise.
func InstrumentingMiddleware(duration metrics.Histogram) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {

			defer func(begin time.Time) {
				duration.With("success", fmt.Sprint(err == nil)).Observe(time.Since(begin).Seconds())
			}(time.Now())
			return next(ctx, request)
		}
	}
}

// LoggingMiddleware returns an endpoint middleware that logs the
// duration of each invocation, and the resulting error, if any.
func LoggingMiddleware(name string, logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(begin time.Time) {
				logger.Log("name", name, "transport_error", err, "took", time.Since(begin))
			}(time.Now())
			return next(ctx, request)
		}
	}
}

func MiscMiddleware(name string, logger log.Logger, duration metrics.Histogram) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		next = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 1))(next)
		next = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(next)
		//next = opentracing.TraceServer(otTracer, name)(next)
		//next = zipkin.TraceEndpoint(zipkinTracer, name)(next)
		next = LoggingMiddleware(name, logger)(next)
		if duration != nil {
			next = InstrumentingMiddleware(duration)(next)
		}
		return next
	}
}

//type UserSvcMiddleware func(svc user.Service) user.Service
