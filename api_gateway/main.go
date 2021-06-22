package main

import (
	"net/http"
	"os"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	zipkin "github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"

	apihttp "gokit-ddd-demo/api_gateway/api/http"
	"gokit-ddd-demo/api_gateway/svc"
	"gokit-ddd-demo/lib/kitx"
	ordergrpc "gokit-ddd-demo/order_svc/api/grpc"
	"gokit-ddd-demo/order_svc/svc/order"
	usergrpc "gokit-ddd-demo/user_svc/api/grpc"
	"gokit-ddd-demo/user_svc/svc/user"
)

func main() {
	errc := make(chan error)

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	// init zipkin trancer
	var tracer *zipkin.Tracer
	{
		zipkinUrl := "http://127.0.0.1:9411/api/v2/spans"
		zipkinEndpoint, err := zipkin.NewEndpoint("api-gateway", "")
		if err != nil {
			panic(err)
		}
		reporter := zipkinhttp.NewReporter(zipkinUrl)
		zipkinTracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(zipkinEndpoint))
		if err != nil {
			panic(err)
		}
		tracer = zipkinTracer
	}

	httpAddr := ":1323"

	// contructor our service
	var (
		usersvc  user.Service
		ordersvc order.Service
	)

	cliOpts := kitx.NewClientOptions(kitx.WithLogger(logger), kitx.WithLoadBalance(3, 5*time.Second), kitx.WithZipkinTracer(tracer))
	{
		instance := []string{":8082"}
		usersvc = usergrpc.NewClient(sd.FixedInstancer(instance), cliOpts)
	}
	{
		instance := []string{":8092"}
		ordersvc = ordergrpc.NewClient(sd.FixedInstancer(instance), cliOpts)
	}
	svc := svc.NewService(usersvc, ordersvc)

	go func() {
		srvOpts := kitx.NewServerOptions(kitx.WithLogger(logger), kitx.WithRateLimit(nil), kitx.WithCircuitBreaker(0), kitx.WithMetrics(nil), kitx.WithZipkinTracer(tracer))
		logger.Log("transport", "HTTP", "addr", httpAddr)
		handler := apihttp.NewHTTPHandler(svc, srvOpts)
		errc <- http.ListenAndServe(httpAddr, handler)
	}()

	err := <-errc
	logger.Log(err)
}
