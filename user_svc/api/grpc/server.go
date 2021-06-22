package grpc

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"

	"gokit-ddd-demo/lib/kitx"
	"gokit-ddd-demo/user_svc/api"
	"gokit-ddd-demo/user_svc/api/grpc/pb"
	"gokit-ddd-demo/user_svc/svc/user"
)

var _ pb.UserSvcServer = (*grpcServer)(nil)

type grpcServer struct {
	find grpctransport.Handler
	get  grpctransport.Handler
}

func makeFindHandler(svc user.Service, opts ...kitx.Option) grpctransport.Handler {
	ep := kitx.ServerEndpoint(func() (endpoint.Endpoint, string) {
		ep := api.MakeFindEndpoint(svc)
		return ep, "user_svc.Find"
	}, opts...)

	return grpctransport.NewServer(ep, decodeFindRequest, encodeFindResponse)
}

func makeGetHandler(svc user.Service, opts ...kitx.Option) grpctransport.Handler {
	ep := kitx.ServerEndpoint(func() (endpoint.Endpoint, string) {
		ep := api.MakeGetEndpoint(svc)
		return ep, "user_svc.Get"
	}, opts...)

	return grpctransport.NewServer(ep, decodeGetRequest, encodeGetResponse)
}

func NewGRPCServer(svc user.Service, opts ...kitx.Option) pb.UserSvcServer {
	srv := &grpcServer{}

	srv.find = makeFindHandler(svc, opts...)
	srv.get = makeGetHandler(svc, opts...)

	return srv
}

func (s *grpcServer) Find(ctx context.Context, req *pb.FindReq) (*pb.FindReply, error) {
	println("server find")
	_, rep, err := s.find.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.FindReply), nil
}

func (s *grpcServer) Get(ctx context.Context, req *pb.ID) (*pb.GetReply, error) {
	println("server get")
	_, rep, err := s.get.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.GetReply), nil
}
