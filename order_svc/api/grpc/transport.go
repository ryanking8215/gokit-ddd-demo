package grpc

import (
	"context"

	"gokit-ddd-demo/lib"
	"gokit-ddd-demo/order_svc/api"
	"gokit-ddd-demo/order_svc/api/grpc/pb"
	"gokit-ddd-demo/order_svc/svc/order"
)

func decodeFindRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	if grpcReq == nil {
		return nil, nil
	}
	id := grpcReq.(*pb.ID)
	return id.GetId(), nil
}

func encodeFindResponse(_ context.Context, response interface{}) (interface{}, error) {
	orders := response.([]*order.Order)
	rsp := &pb.FindReply{}
	for _, o := range orders {
		pbOrder := &pb.Order{Id: o.ID, Userid: o.UserID, Product: o.Product}
		rsp.Order = append(rsp.Order, pbOrder)
	}
	return rsp, nil
}

func decodeGetRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.ID)
	return req.Id, nil
}

func encodeGetResponse(_ context.Context, response interface{}) (interface{}, error) {
	r := response.(*api.Response)
	rsp := &pb.GetReply{}

	if r.Error != nil {
		if err, ok := r.Error.(lib.Error); ok {
			rsp.Err = &pb.Error{Code: int32(err.Code), Reason: err.Message}
		} else {
			rsp.Err = &pb.Error{Code: int32(lib.ErrInternal), Reason: r.Error.Error()}
		}
		return rsp, nil
	}

	o := r.Value.(*order.Order)
	rsp.Order = &pb.Order{Id: o.ID, Userid: o.UserID, Product: o.Product}
	return rsp, nil
}

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
	err := reply.GetErr()
	if err != nil {
		return nil, lib.NewError(int(err.Code), err.Reason)
	}

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
