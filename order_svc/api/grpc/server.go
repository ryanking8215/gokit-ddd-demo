package grpc

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"

	"gokit-ddd-demo/lib/kitx"
	"gokit-ddd-demo/order_svc/api"
	"gokit-ddd-demo/order_svc/api/grpc/pb"
	"gokit-ddd-demo/order_svc/svc/order"
)

var _ pb.OrderSvcServer = (*grpcServer)(nil)

type grpcServer struct {
	find grpctransport.Handler
	get  grpctransport.Handler
}

func makeFindHandler(svc order.Service, opts ...kitx.Option) grpctransport.Handler {
	ep := kitx.ServerEndpoint(func() (endpoint.Endpoint, string) {
		ep := api.MakeFindEndpoint(svc)
		return ep, "order_svc.Find"
	}, opts...)

	return grpctransport.NewServer(ep, decodeFindRequest, encodeFindResponse)
}

func makeGetHandler(svc order.Service, opts ...kitx.Option) grpctransport.Handler {
	ep := kitx.ServerEndpoint(func() (endpoint.Endpoint, string) {
		ep := api.MakeGetEndpoint(svc)
		return ep, "order_svc.Get"
	}, opts...)

	return grpctransport.NewServer(ep, decodeGetRequest, encodeGetResponse)
}

func NewGRPCServer(svc order.Service, opts ...kitx.Option) pb.OrderSvcServer {
	srv := &grpcServer{}

	srv.find = makeFindHandler(svc, opts...)
	srv.get = makeGetHandler(svc, opts...)

	return srv
}

func (s *grpcServer) Find(ctx context.Context, req *pb.ID) (*pb.FindReply, error) {
	_, rep, err := s.find.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.FindReply), nil
}

func (s *grpcServer) Get(ctx context.Context, req *pb.ID) (*pb.GetReply, error) {
	_, rep, err := s.get.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.GetReply), nil
}
