package grpc

import (
	"context"

	"gokit-ddd-demo/user_svc/api/grpc/pb"
	//"gokit-ddd-demo/user_svc/domain/user"
	"gokit-ddd-demo/user_svc/domain/models"
)

func decodeFindRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	return nil, nil
}

func encodeFindResponse(_ context.Context, response interface{}) (interface{}, error) {
	users := response.([]*models.User)
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
	u := response.(*models.User)
	rsp := &pb.GetReply{}
	rsp.User = &pb.User{Id: u.ID, Name: u.Name}
	return rsp, nil
}
