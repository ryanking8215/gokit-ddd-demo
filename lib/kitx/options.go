package kitx

import (
	"time"

	"github.com/go-kit/kit/log"

	stdopentracing "github.com/opentracing/opentracing-go"
)

type Option interface {
	apply(*options)
}

type options struct {
	circuitBreakerOption
	rateLimitOption
	openTracingOption
	loggerOption
	loadBalanceOption
	metricsOption
}

type circuitBreakerOption struct {
	enable  bool
	timeout time.Duration
}

func (o circuitBreakerOption) apply(opts *options) {
	opts.circuitBreakerOption = o
}

func WithCircuitBreaker(timeout time.Duration) Option {
	return circuitBreakerOption{true, timeout}
}

type rateLimitOption struct {
	enable   bool
	every    time.Duration
	tokenCnt int
}

func (o rateLimitOption) apply(opts *options) {
	opts.rateLimitOption = o
}

func WithRateLimit(every time.Duration, tokenCnt int) Option {
	return rateLimitOption{true, every, tokenCnt}
}

type openTracingOption struct {
	otTracer stdopentracing.Tracer
}

func (o openTracingOption) apply(opts *options) {
	opts.openTracingOption = o
}

func WithOpenTracing(otTracer stdopentracing.Tracer) Option {
	return openTracingOption{otTracer}
}

// func WithZipkin() Option {
// }

type loggerOption struct {
	logger log.Logger
}

func (o loggerOption) apply(opts *options) {
	opts.loggerOption = o
}

func WithLogger(logger log.Logger) Option {
	return loggerOption{logger}
}

type loadBalanceOption struct {
	enable       bool
	retryMax     int
	retryTimeout time.Duration
}

func (o loadBalanceOption) apply(opts *options) {
	opts.loadBalanceOption = o
}

func WithLoadBalance(retryMax int, retryTimeout time.Duration) Option {
	return loadBalanceOption{true, retryMax, retryTimeout}
}

type metricsOption struct {
	enable bool
}

func (o metricsOption) apply(opts *options) {
	opts.metricsOption = o
}

func WithMetrics() Option {
	return metricsOption{true}
}
