package api

import "github.com/labstack/echo/v4"

func SetupRouter(e *echo.Echo) {
	api := e.Group("/api")
	api.GET("/users", GetUsers)
	api.GET("/users/:id", GetUser)
	api.GET("/users/:id/orders", GetUserOrders)
	api.GET("/users/:id/orders/:orderid", GetUserOrder)
}
