package main

import (
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"

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

	httpAddr := ":1323"

	var (
		usersvc  user.Service
		ordersvc order.Service
	)

	{
		instance := []string{":8082"}
		usersvc = usergrpc.NewClient(sd.FixedInstancer(instance), kitx.WithLogger(logger))
	}
	{
		instance := []string{":8092"}
		ordersvc = ordergrpc.NewClient(sd.FixedInstancer(instance), kitx.WithLogger(logger))
	}
	svc := svc.NewService(usersvc, ordersvc)

	go func() {
		logger.Log("transport", "HTTP", "addr", httpAddr)
		handler := apihttp.NewHTTPHandler(svc)
		errc <- http.ListenAndServe(httpAddr, handler)
	}()

	err := <-errc
	logger.Log(err)
}
