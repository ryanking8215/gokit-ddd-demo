package grpc

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/tracing/zipkin"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	grpcpool "github.com/processout/grpc-go-pool"
	"google.golang.org/grpc/metadata"

	"gokit-ddd-demo/lib/kitx"
	orderpb "gokit-ddd-demo/order_svc/api/grpc/pb"
	"gokit-ddd-demo/order_svc/svc/order"
)

var _ order.Service = (*Client)(nil)

type Client struct {
	find endpoint.Endpoint
	get  endpoint.Endpoint
}

func (c *Client) Find(ctx context.Context, userID int64) ([]*order.Order, error) {
	rsp, err := c.find(ctx, userID)
	if err != nil {
		return nil, err
	}
	return rsp.([]*order.Order), err
}

func (c *Client) Get(ctx context.Context, id int64) (*order.Order, error) {
	rsp, err := c.get(ctx, id)
	if err != nil {
		return nil, err
	}
	return rsp.(*order.Order), err
}

func (c *Client) Delete(ctx context.Context, id int64) error {
	return nil
}

func (c *Client) Save(ctx context.Context, o *order.Order) error {
	return nil
}

func NewClient(instancer sd.Instancer, opts *kitx.ClientOptions) *Client {
	c := &Client{}

	var options []grpctransport.ClientOption
	tracer := opts.ZipkinTracer()
	if tracer != nil {
		options = append(options, zipkin.GRPCClientTrace(tracer))
		options = append(options, grpctransport.ClientFinalizer(kitx.GRPCClientFinalizer)) // take care of the conn from grpc conn pool
	}

	c.find = kitx.GRPCClientEndpoint(instancer, func(conn *grpcpool.ClientConn) (endpoint.Endpoint, string) {
		return grpctransport.NewClient(
			conn.ClientConn,
			"pb.OrderSvc",
			"Find",
			encodeFindRequest,
			decodeFindResponse,
			orderpb.FindReply{},
			append(options, grpctransport.ClientBefore(func(ctx context.Context, _ *metadata.MD) context.Context {
				return context.WithValue(ctx, kitx.GRPCConnKey, conn) // inject the conn to ctx
			}))...,
		).Endpoint(), "order_srv.rpc.Find"
	}, opts)

	c.get = kitx.GRPCClientEndpoint(instancer, func(conn *grpcpool.ClientConn) (endpoint.Endpoint, string) {
		return grpctransport.NewClient(
			conn.ClientConn,
			"pb.UserSvc",
			"Get",
			encodeGetRequest,
			decodeGetResponse,
			orderpb.GetReply{},
			append(options, grpctransport.ClientBefore(func(ctx context.Context, _ *metadata.MD) context.Context {
				return context.WithValue(ctx, kitx.GRPCConnKey, conn) // inject the conn to ctx
			}))...,
		).Endpoint(), "order_srv.rpc.Get"
	}, opts)

	return c
}
