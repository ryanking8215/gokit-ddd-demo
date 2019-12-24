package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"gokit-ddd-demo/order_svc/domain/order"
	"gokit-ddd-demo/user_svc/domain/user"
)

var UserSvc user.Service
var OrderSvc order.Service

func GetUsers(ctx echo.Context) error {
	withOrders := false
	str := ctx.QueryParam("with_orders")
	if str == "true" {
		withOrders = true
	}
	users, err := UserSvc.Find(context.TODO())
	if err != nil {
		return NewInternalServerError(err.Error())
	}
	usersRsp := make([]*User, 0, len(users))
	for _, u := range users {
		usersRsp = append(usersRsp, &User{u, nil})
	}

	if withOrders {
		for _, u := range usersRsp {
			orders, err := OrderSvc.Find(context.TODO(), u.ID)
			if err != nil {
				return NewInternalServerError(err.Error())
			}
			u.Orders = orders
		}
	}

	return ctx.JSON(http.StatusOK, usersRsp)
}

func GetUser(ctx echo.Context) error {
	str := ctx.Param("id")
	if str == "" {
		return NewInvalidParamError("id is empty")
	}
	userID, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return NewInvalidParamError(err.Error())
	}
	user, err := UserSvc.Get(context.TODO(), userID)
	if err != nil {
		return NewInternalServerError(err.Error())
	}
	return ctx.JSON(http.StatusOK, user)
}

func GetUserOrders(ctx echo.Context) error {
	str := ctx.Param("id")
	if str == "" {
		return NewInvalidParamError("id is empty")
	}
	userID, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return NewInvalidParamError(err.Error())
	}
	_ = userID
	// TODO
	return nil
}

func GetUserOrder(ctx echo.Context) error {
	// TODO
	return nil
}
