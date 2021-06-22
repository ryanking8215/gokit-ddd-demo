package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"text/tabwriter"

	"github.com/go-kit/kit/log"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"google.golang.org/grpc"

	"gokit-ddd-demo/lib/kitx"
	apigrpc "gokit-ddd-demo/user_svc/api/grpc"
	"gokit-ddd-demo/user_svc/api/grpc/pb"
	apihttp "gokit-ddd-demo/user_svc/api/http"
	"gokit-ddd-demo/user_svc/infras/repo/inmem"
	"gokit-ddd-demo/user_svc/svc/user"
)

func main() {
	errc := make(chan error)

	fs := flag.NewFlagSet("usersvc", flag.ExitOnError)
	var (
		//debugAddr      = fs.String("debug.addr", ":8080", "Debug and metrics listen address")
		httpAddr = fs.String("http-addr", ":8081", "HTTP listen address")
		grpcAddr = fs.String("grpc-addr", ":8082", "gRPC listen address")
		//jsonRPCAddr    = fs.String("jsonrpc-addr", ":8084", "JSON RPC listen address")
		//thriftProtocol = fs.String("thrift-protocol", "binary", "binary, compact, json, simplejson")
		//thriftBuffer   = fs.Int("thrift-buffer", 0, "0 for unbuffered")
		//thriftFramed   = fs.Bool("thrift-framed", false, "true to enable framing")
		//zipkinV2URL    = fs.String("zipkin-url", "", "Enable Zipkin v2 tracing (zipkin-go) using a Reporter URL e.g. http://localhost:9411/api/v2/spans")
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

	// init zipkin trancer
	var tracer *zipkin.Tracer
	{
		zipkinUrl := "http://127.0.0.1:9411/api/v2/spans"
		zipkinEndpoint, err := zipkin.NewEndpoint("user-svc", "")
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

	// init infras and svc
	var (
		userRepo = inmem.NewUserRepo()
		//orderGRPCClient = ordergrpc.NewGRPCClient(sd.FixedInstancer{":8092"}, 3, 5*time.Second, logger)
		userSvc = user.NewUserService(userRepo)
	)
	_ = userRepo.InitMockData([]*user.User{
		{ID: 1, Name: "user1"},
		{ID: 2, Name: "user2"},
	})

	srvOpts := kitx.NewServerOptions(kitx.WithLogger(logger), kitx.WithRateLimit(nil), kitx.WithCircuitBreaker(0), kitx.WithMetrics(nil), kitx.WithZipkinTracer(tracer))

	var (
		grpcServer  = apigrpc.NewGRPCServer(userSvc, srvOpts)
		httpHandler = apihttp.NewHTTPHandler(userSvc, srvOpts)
	)

	// The HTTP listener mounts the Go kit HTTP handler we created.
	go func() {
		httpListener, err := net.Listen("tcp", *httpAddr)
		if err != nil {
			errc <- err
			return
		}
		errc <- http.Serve(httpListener, httpHandler)
	}()

	// The gRPC listener mounts the Go kit gRPC server we created.
	go func() {
		grpcListener, err := net.Listen("tcp", *grpcAddr)
		if err != nil {
			errc <- err
			return
		}

		//baseServer := grpc.NewServer(grpc.UnaryInterceptor(kitgrpc.Interceptor))
		s := grpc.NewServer()
		pb.RegisterUserSvcServer(s, grpcServer)
		errc <- s.Serve(grpcListener)
	}()

	logger.Log("during", "run")

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
