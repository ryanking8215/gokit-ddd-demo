package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"text/tabwriter"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"google.golang.org/grpc"

	"gokit-ddd-demo/lib/kitx"
	apigrpc "gokit-ddd-demo/order_svc/api/grpc"
	"gokit-ddd-demo/order_svc/api/grpc/pb"
	apihttp "gokit-ddd-demo/order_svc/api/http"
	"gokit-ddd-demo/order_svc/infras/repo/inmem"
	"gokit-ddd-demo/order_svc/svc/order"
)

func main() {
	fs := flag.NewFlagSet("ordersvc", flag.ExitOnError)
	var (
		//debugAddr      = fs.String("debug.addr", ":8080", "Debug and metrics listen address")
		httpAddr = fs.String("http-addr", ":8091", "HTTP listen address")
		grpcAddr = fs.String("grpc-addr", ":8092", "gRPC listen address")
		//jsonRPCAddr    = fs.String("jsonrpc-addr", ":8084", "JSON RPC listen address")
		//thriftProtocol = fs.String("thrift-protocol", "binary", "binary, compact, json, simplejson")
		//thriftBuffer   = fs.Int("thrift-buffer", 0, "0 for unbuffered")
		//thriftFramed   = fs.Bool("thrift-framed", false, "true to enable framing")
		zipkinV2URL = fs.String("zipkin-url", "", "Enable Zipkin v2 tracing (zipkin-go) using a Reporter URL e.g. http://localhost:9411/api/v2/spans")
		//zipkinV1URL    = fs.String("zipkin-v1-url", "", "Enable Zipkin v1 tracing (zipkin-go-opentracing) using a collector URL e.g. http://localhost:9411/api/v1/spans")
		//lightstepToken = fs.String("lightstep-token", "", "Enable LightStep tracing via a LightStep access token")
		//appdashAddr    = fs.String("appdash-addr", "", "Enable Appdash tracing via an Appdash server host:port")
	)
	fs.Usage = usageFor(fs, os.Args[0]+" [flags]")
	fs.Parse(os.Args[1:])

	// Create a single logger, which we'll use and give to other components.
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var tracer *zipkin.Tracer
	{
		if *zipkinV2URL != "" {
			zipkinEndpoint, err := zipkin.NewEndpoint("order-svc", "")
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

	// init infras and domain
	var (
		orderRepo = inmem.NewOrderRepo()
		orderSvc  = order.NewOrderService(orderRepo)
	)
	orderRepo.InitMockData([]*order.Order{
		{ID: 1, UserID: 1, Product: "product1"},
		{ID: 2, UserID: 1, Product: "product2"},
		{ID: 3, UserID: 2, Product: "product1"},
	})

	srvOpts := kitx.NewServerOptions(kitx.WithLogger(logger), kitx.WithRateLimit(nil), kitx.WithCircuitBreaker(0), kitx.WithMetrics(nil), kitx.WithZipkinTracer(tracer))

	var (
		grpcServer  = apigrpc.NewGRPCServer(orderSvc, srvOpts)
		httpHandler = apihttp.NewHTTPHandler(orderSvc, srvOpts)
	)

	{
		// The HTTP listener mounts the Go kit HTTP handler we created.
		httpListener, err := net.Listen("tcp", *httpAddr)
		if err != nil {
			logger.Log("transport", "HTTP", "during", "Listen", "err", err)
			os.Exit(1)
		}
		go func() {
			if err := http.Serve(httpListener, httpHandler); err != nil {
				logger.Log("transport", "http", "during", "serve", "err", err)
				os.Exit(1)
			}
		}()
	}
	{
		// The gRPC listener mounts the Go kit gRPC server we created.
		grpcListener, err := net.Listen("tcp", *grpcAddr)
		if err != nil {
			logger.Log("transport", "gRPC", "during", "Listen", "err", err)
			os.Exit(1)
		}
		//baseServer := grpc.NewServer(grpc.UnaryInterceptor(kitgrpc.Interceptor))
		s := grpc.NewServer()
		pb.RegisterOrderSvcServer(s, grpcServer)
		go func() {
			if err := s.Serve(grpcListener); err != nil {
				logger.Log("transport", "gRPC", "during", "serve", "err", err)
				os.Exit(1)
			}
		}()
	}

	logger.Log("during", "run")

	for {
		time.Sleep(10 * time.Second)
	}
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
