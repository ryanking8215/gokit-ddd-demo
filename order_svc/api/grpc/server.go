package grpc

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/tracing/zipkin"
	"github.com/go-kit/kit/transport"
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

func makeFindHandler(svc order.Service, options []grpctransport.ServerOption, opts *kitx.ServerOptions) grpctransport.Handler {
	ep := kitx.ServerEndpoint(func() (endpoint.Endpoint, string) {
		ep := api.MakeFindEndpoint(svc)
		return ep, "order_svc.Find"
	}, opts)

	return grpctransport.NewServer(ep, decodeFindRequest, encodeFindResponse, options...)
}

func makeGetHandler(svc order.Service, options []grpctransport.ServerOption, opts *kitx.ServerOptions) grpctransport.Handler {
	ep := kitx.ServerEndpoint(func() (endpoint.Endpoint, string) {
		ep := api.MakeGetEndpoint(svc)
		return ep, "order_svc.Get"
	}, opts)

	return grpctransport.NewServer(ep, decodeGetRequest, encodeGetResponse, options...)
}

func NewGRPCServer(svc order.Service, opts *kitx.ServerOptions) pb.OrderSvcServer {
	srv := &grpcServer{}

	logger := opts.Logger()
	tracer := opts.ZipkinTracer()

	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	if tracer != nil {
		options = append(options, zipkin.GRPCServerTrace(tracer))
	}

	srv.find = makeFindHandler(svc, options, opts)
	srv.get = makeGetHandler(svc, options, opts)

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
