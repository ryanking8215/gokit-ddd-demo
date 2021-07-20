package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"text/tabwriter"
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
	fs := flag.NewFlagSet("ordersvc", flag.ExitOnError)
	var (
		//debugAddr      = fs.String("debug.addr", ":8080", "Debug and metrics listen address")
		httpAddr = fs.String("http-addr", ":1323", "HTTP listen address")
		//jsonRPCAddr    = fs.String("jsonrpc-addr", ":8084", "JSON RPC listen address")
		//thriftProtocol = fs.String("thrift-protocol", "binary", "binary, compact, json, simplejson")
		//thriftBuffer   = fs.Int("thrift-buffer", 0, "0 for unbuffered")
		//thriftFramed   = fs.Bool("thrift-framed", false, "true to enable framing")
		zipkinV2URL      = fs.String("zipkin-url", "", "Enable Zipkin v2 tracing (zipkin-go) using a Reporter URL e.g. http://localhost:9411/api/v2/spans")
		usersvcInstance  = fs.String("usersvc", ":8082", "Instance of user service")
		ordersvcInstance = fs.String("ordersvc", ":8092", "Instance of order service")
		//zipkinV1URL    = fs.String("zipkin-v1-url", "", "Enable Zipkin v1 tracing (zipkin-go-opentracing) using a collector URL e.g. http://localhost:9411/api/v1/spans")
		//lightstepToken = fs.String("lightstep-token", "", "Enable LightStep tracing via a LightStep access token")
		//appdashAddr    = fs.String("appdash-addr", "", "Enable Appdash tracing via an Appdash server host:port")
	)
	fs.Usage = usageFor(fs, os.Args[0]+" [flags]")
	fs.Parse(os.Args[1:])
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
		if *zipkinV2URL != "" {
			zipkinEndpoint, err := zipkin.NewEndpoint("api-gateway", "")
			if err != nil {
				panic(err)
			}
			reporter := zipkinhttp.NewReporter(*zipkinV2URL)
			zipkinTracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(zipkinEndpoint))
			if err != nil {
				panic(err)
			}
			tracer = zipkinTracer
		} else {
			zipkinTracer, _ := zipkin.NewTracer(nil, zipkin.WithNoopTracer(true))
			tracer = zipkinTracer
		}
	}

	// contructor our service
	var (
		usersvc  user.Service
		ordersvc order.Service
	)

	cliOpts := kitx.NewClientOptions(kitx.WithLogger(logger), kitx.WithLoadBalance(3, 5*time.Second), kitx.WithZipkinTracer(tracer))
	{
		instance := []string{*usersvcInstance}
		usersvc = usergrpc.NewClient(sd.FixedInstancer(instance), cliOpts)
	}
	{
		instance := []string{*ordersvcInstance}
		ordersvc = ordergrpc.NewClient(sd.FixedInstancer(instance), cliOpts)
	}
	svc := svc.NewService(usersvc, ordersvc)

	go func() {
		srvOpts := kitx.NewServerOptions(kitx.WithLogger(logger), kitx.WithRateLimit(nil), kitx.WithCircuitBreaker(0), kitx.WithMetrics(nil), kitx.WithZipkinTracer(tracer))
		logger.Log("transport", "HTTP", "addr", *httpAddr)
		handler := apihttp.NewHTTPHandler(svc, srvOpts)
		errc <- http.ListenAndServe(*httpAddr, handler)
	}()

	err := <-errc
	logger.Log(err)
}

func usageFor(fs *flag.FlagSet, short string) func() {
	return func() {
		fmt.Fprintf(os.Stderr, "USAGE\n")
		fmt.Fprintf(os.Stderr, "  %s\n", short)
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "FLAGS\n")
		w := tabwriter.NewWriter(os.Stderr, 0, 2, 2, ' ', 0)
		fs.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(w, "\t-%s %s\t%s\n", f.Name, f.DefValue, f.Usage)
		})
		w.Flush()
		fmt.Fprintf(os.Stderr, "\n")
	}
}
