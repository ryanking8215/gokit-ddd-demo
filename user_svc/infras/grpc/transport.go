package grpc

import (
	"context"
	pb "gokit-ddd-demo/user_svc/api/grpc/pb"
	"gokit-ddd-demo/user_svc/domain/user"
)

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
