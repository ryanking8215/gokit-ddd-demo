package grpc

import (
	"context"
	"io"
	"time"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/sony/gobreaker"
	"google.golang.org/grpc"

	pb "gokit-ddd-demo/order_svc/api/grpc/pb"
	"gokit-ddd-demo/order_svc/domain/order"
	//"gokit-ddd-demo/user_svc/domain/user"
)

var _ order.Service = (*grpcClient)(nil)

type grpcClient struct {
	retryMax     int
	retryTimeout time.Duration
	findEndpoint endpoint.Endpoint
	getEndpoint  endpoint.Endpoint
}

//func NewGRPCClientSimple(instance string, logger log.Logger) *grpcClient {
//	c := &grpcClient{}
//	if instance != "" {
//		c.retryMax = 2
//		c.retryTimeout = 5 * time.Second
//		instancer := sd.FixedInstancer{instance}
//		c.initEndpoints("", instancer, logger)
//	}
//	return c
//}

func NewGRPCClient(instancer sd.Instancer, retryMax int, retryTimeout time.Duration, logger log.Logger) *grpcClient {
	c := &grpcClient{}
	c.retryMax = retryMax
	c.retryTimeout = retryTimeout
	c.initEndpoints(instancer, logger)
	return c
}

func (c *grpcClient) initEndpoints(instancer sd.Instancer, logger log.Logger) error {
	{
		factory := newFactory(func(conn *grpc.ClientConn) endpoint.Endpoint {
			// build endpoint
			ep := grpctransport.NewClient(
				conn,
				"pb.OrderSvc",
				"Find",
				encodeFindRequest,
				decodeFindResponse,
				pb.FindReply{},
				//append(options, grpctransport.ClientBefore(opentracing.ContextToGRPC(otTracer, logger)))...,
			).Endpoint()
			//ep = opentracing.TraceClient(otTracer, "Sum")(sumEndpoint)
			//ep = limiter(sumEndpoint)
			ep = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
				Name:    "ordersvc.Find",
				Timeout: 30 * time.Second,
			}))(ep)
			return ep
		})
		//if instance != "" {
		//	ep := simpleEndpoint(instance, factory)
		//	c.findEndpoint = ep
		//}
		if instancer != nil {
			endpointer := sd.NewEndpointer(instancer, factory, logger)
			balancer := lb.NewRoundRobin(endpointer)
			retry := lb.Retry(c.retryMax, c.retryTimeout, balancer)
			c.findEndpoint = retry
		}
	}
	{
		factory := newFactory(func(conn *grpc.ClientConn) endpoint.Endpoint {
			// build endpoint
			ep := grpctransport.NewClient(
				conn,
				"pb.OrderSvc",
				"Get",
				encodeGetRequest,
				decodeGetResponse,
				pb.GetReply{},
				//append(options, grpctransport.ClientBefore(opentracing.ContextToGRPC(otTracer, logger)))...,
			).Endpoint()
			//ep = opentracing.TraceClient(otTracer, "Sum")(sumEndpoint)
			//ep = limiter(sumEndpoint)
			ep = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
				Name:    "usersvc.Get",
				Timeout: 30 * time.Second,
			}))(ep)
			return ep
		})
		//if instance != "" {
		//	ep := simpleEndpoint(instance, factory)
		//	c.getEndpoint = ep
		//}
		if instancer != nil {
			endpointer := sd.NewEndpointer(instancer, factory, logger)
			balancer := lb.NewRoundRobin(endpointer)
			retry := lb.Retry(c.retryMax, c.retryTimeout, balancer)
			c.getEndpoint = retry
		}
	}
	return nil
}

func (c *grpcClient) Find(ctx context.Context, userID int64) ([]*order.Order, error) {
	rsp, err := c.findEndpoint(ctx, userID)
	if err != nil {
		return nil, err
	}
	return rsp.([]*order.Order), err
}

func (c *grpcClient) Get(ctx context.Context, id int64) (*order.Order, error) {
	rsp, err := c.getEndpoint(ctx, id)
	if err != nil {
		return nil, err
	}
	return rsp.(*order.Order), err
}

func (c *grpcClient) Delete(ctx context.Context, id int64) error {
	return nil
}

func (c *grpcClient) Save(ctx context.Context, o *order.Order) error {
	return nil
}

func newFactory(makeEndpoint func(conn *grpc.ClientConn) endpoint.Endpoint) sd.Factory {
	return func(instance string) (i endpoint.Endpoint, closer io.Closer, e error) {
		conn, err := grpc.Dial(instance, grpc.WithInsecure())
		if err != nil {
			return nil, nil, err
		}
		ep := makeEndpoint(conn)
		return ep, conn, nil
	}
}

//func simpleEndpoint(instance string, factory sd.Factory) endpoint.Endpoint {
//	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
//		ep, closer, err := factory(instance)
//		if err != nil {
//			return nil, err
//		}
//		defer closer.Close()
//		return ep(ctx, request)
//	}
//}
