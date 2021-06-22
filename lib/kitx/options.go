package kitx

import (
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/ratelimit"
	"github.com/openzipkin/zipkin-go"

	stdopentracing "github.com/opentracing/opentracing-go"
)

type Option interface {
	apply(*options)
}

type options struct {
	circuitBreakerOption
	rateLimitOption
	openTracingOption
	zipkinTracerOption
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
	limiter ratelimit.Allower
}

func (o rateLimitOption) apply(opts *options) {
	opts.rateLimitOption = o
}

func WithRateLimit(limiter ratelimit.Allower) Option {
	return rateLimitOption{limiter}
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
	retryMax int
	timeout  time.Duration
}

func (o loadBalanceOption) apply(opts *options) {
	opts.loadBalanceOption = o
}

func WithLoadBalance(retryMax int, timeout time.Duration) Option {
	return loadBalanceOption{retryMax, timeout}
}

type metricsOption struct {
	histogram metrics.Histogram
}

func (o metricsOption) apply(opts *options) {
	opts.metricsOption = o
}

func WithMetrics(histogram metrics.Histogram) Option {
	return metricsOption{histogram}
}

type zipkinTracerOption struct {
	tracer *zipkin.Tracer
}

func (o zipkinTracerOption) apply(opts *options) {
	opts.zipkinTracerOption = o
}

func WithZipkinTracer(tracer *zipkin.Tracer) Option {
	return zipkinTracerOption{tracer}
}

type ServerOptions struct {
	options
}

func NewServerOptions(opts ...Option) *ServerOptions {
	so := &ServerOptions{}
	for _, o := range opts {
		o.apply(&so.options)
	}
	return so
}

func (o *ServerOptions) Logger() log.Logger {
	return o.loggerOption.logger
}

func (o *ServerOptions) ZipkinTracer() *zipkin.Tracer {
	return o.zipkinTracerOption.tracer
}

type ClientOptions struct {
	options
}

func NewClientOptions(opts ...Option) *ClientOptions {
	co := &ClientOptions{}
	for _, o := range opts {
		o.apply(&co.options)
	}
	return co
}

func (o *ClientOptions) Logger() log.Logger {
	return o.loggerOption.logger
}

func (o *ClientOptions) ZipkinTracer() *zipkin.Tracer {
	return o.zipkinTracerOption.tracer
}
