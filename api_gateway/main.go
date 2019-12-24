package main

import (
	"os"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"gokit-ddd-demo/api_gateway/api"
	ordergrpc "gokit-ddd-demo/order_svc/infras/grpc"
	usergrpc "gokit-ddd-demo/user_svc/infras/grpc"
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	//e.Use(middleware.Recover())
	e.HTTPErrorHandler = api.ErrorHandle

	// Routes
	api.SetupRouter(e)

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}
	setupService(logger)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

func setupService(logger log.Logger) {
	{
		instance := ":8082"
		api.UserSvc = usergrpc.NewGRPCClient(sd.FixedInstancer{instance}, 3, 5*time.Second, logger)
	}
	{
		instance := ":8092"
		api.OrderSvc = ordergrpc.NewGRPCClient(sd.FixedInstancer{instance}, 3, 5*time.Second, logger)
	}
}
