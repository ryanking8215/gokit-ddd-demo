package grpc

import (
	"context"

	"gokit-ddd-demo/user_svc/api/grpc/pb"
	"gokit-ddd-demo/user_svc/domain/user"
)

func decodeFindRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.FindReq)
	return req, nil
}

func encodeFindResponse(_ context.Context, response interface{}) (interface{}, error) {
	users := response.([]*user.User)
	rsp := &pb.FindReply{}
	for _, u := range users {
		pbUser := &pb.User{Id: u.ID, Name: u.Name}
		rsp.Users = append(rsp.Users, pbUser)
	}
	return rsp, nil
}

func decodeGetRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.ID)
	return req.Id, nil
}

func encodeGetResponse(_ context.Context, response interface{}) (interface{}, error) {
	u := response.(*user.User)
	rsp := &pb.GetReply{}
	rsp.User = &pb.User{Id: u.ID, Name: u.Name}
	return rsp, nil
}
