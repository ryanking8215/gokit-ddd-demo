package grpc

import (
	"context"

	pb "gokit-ddd-demo/order_svc/api/grpc/pb"
	//"gokit-ddd-demo/order_svc/domain/models"
	"gokit-ddd-demo/order_svc/domain/order"
)

func encodeFindRequest(_ context.Context, request interface{}) (interface{}, error) {
	userID := request.(int64)
	return &pb.ID{Id: userID}, nil
}

func decodeFindResponse(_ context.Context, response interface{}) (interface{}, error) {
	reply := response.(*pb.FindReply)
	var orders []*order.Order
	for _, o := range reply.GetOrder() {
		orders = append(orders, &order.Order{
			ID:      o.Id,
			Product: o.Product,
		})
	}
	return orders, nil
}

func encodeGetRequest(_ context.Context, request interface{}) (interface{}, error) {
	id := request.(int64)
	return &pb.ID{Id: id}, nil
}

func decodeGetResponse(_ context.Context, response interface{}) (interface{}, error) {
	reply := response.(*pb.GetReply)
	o := reply.GetOrder()
	var item *order.Order
	if o != nil {
		item = &order.Order{
			ID:      o.Id,
			Product: o.Product,
		}
	}
	return item, nil
}
