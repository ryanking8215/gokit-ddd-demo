package grpc

import (
	"context"

	"gokit-ddd-demo/user_svc/api/grpc/pb"
	"gokit-ddd-demo/user_svc/svc/user"
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

func encodeFindRequest(_ context.Context, request interface{}) (interface{}, error) {
	return nil, nil
}

func decodeFindResponse(_ context.Context, response interface{}) (interface{}, error) {
	reply := response.(*pb.FindReply)
	var users []*user.User
	for _, u := range reply.GetUsers() {
		user := &user.User{
			ID:   u.Id,
			Name: u.Name,
		}
		users = append(users, user)
	}
	return users, nil
}

func encodeGetRequest(_ context.Context, request interface{}) (interface{}, error) {
	id := request.(int64)
	return &pb.ID{Id: id}, nil
}

func decodeGetResponse(_ context.Context, response interface{}) (interface{}, error) {
	reply := response.(*pb.GetReply)
	u := reply.GetUser()
	var item *user.User
	if u != nil {
		item = &user.User{
			ID:   u.Id,
			Name: u.Name,
		}
	}
	return item, nil
}
