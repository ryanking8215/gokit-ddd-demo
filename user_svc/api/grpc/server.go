package grpc

import (
	"context"
	"log"

	//"github.com/go-kit/kit/tracing/opentracing"
	//"github.com/go-kit/kit/tracing/zipkin"
	//"github.com/go-kit/kit/transport"
	kitlog "github.com/go-kit/kit/log"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/golang/protobuf/ptypes/empty"

	"gokit-ddd-demo/user_svc/api/endpoint"
	"gokit-ddd-demo/user_svc/api/grpc/pb"
	"gokit-ddd-demo/user_svc/domain/user"
)

var _ pb.UserSvcServer = (*grpcServer)(nil)

type grpcServer struct {
	find grpctransport.Handler
	get  grpctransport.Handler
}

func NewGRPCServer(svc user.Service, logger kitlog.Logger) pb.UserSvcServer {
	// Zipkin GRPC Server Trace can either be instantiated per gRPC method with a
	// provided operation name or a global tracing service can be instantiated
	// without an operation name and fed to each Go kit gRPC server as a
	// ServerOption.
	// In the latter case, the operation name will be the endpoint's grpc method
	// path if used in combination with the Go kit gRPC Interceptor.
	//
	// In this example, we demonstrate a global Zipkin tracing service with
	// Go kit gRPC Interceptor.
	//zipkinServer := zipkin.GRPCServerTrace(zipkinTracer)
	//
	//options := []grpctransport.ServerOption{
	//	grpctransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	//	zipkinServer,
	//}

	srv := &grpcServer{}
	{
		ep := endpoint.MakeFindEndpoint(svc)
		ep = endpoint.MiscMiddleware("usersvc.Find", logger, nil)(ep)
		srv.find = grpctransport.NewServer(
			ep,
			decodeFindRequest,
			encodeFindResponse,
			//append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(otTracer, "Sum", logger)))...,
		)
	}
	{
		ep := endpoint.MakeGetEndpoint(svc)
		ep = endpoint.MiscMiddleware("usersvc.Get", logger, nil)(ep)
		srv.get = grpctransport.NewServer(
			ep,
			decodeGetRequest,
			encodeGetResponse,
			//append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(otTracer, "Sum", logger)))...,
		)
	}
	return srv
}

func (s *grpcServer) Find(ctx context.Context, req *empty.Empty) (*pb.FindReply, error) {
	log.Println(">>>>> find")
	_, rep, err := s.find.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.FindReply), nil
}

func (s *grpcServer) Get(ctx context.Context, req *pb.ID) (*pb.GetReply, error) {
	log.Println(">>>>> get")
	_, rep, err := s.get.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.GetReply), nil
}
