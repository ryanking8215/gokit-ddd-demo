package main

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"os"
	"time"

	"gokit-ddd-demo/api_gateway/api"
	"gokit-ddd-demo/api_gateway/svc/usersvc"
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
		api.UserSvc = usersvc.NewGRPCClient(sd.FixedInstancer{instance}, 3, 5*time.Second, logger)
	}
}
