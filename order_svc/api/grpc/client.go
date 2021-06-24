package grpc

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/tracing/zipkin"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"

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
	}

	c.find = kitx.GRPCClientEndpoint(instancer, func(conn *grpc.ClientConn) (endpoint.Endpoint, string) {
		return grpctransport.NewClient(
			conn,
			"pb.OrderSvc",
			"Find",
			encodeFindRequest,
			decodeFindResponse,
			orderpb.FindReply{},
			options...,
		).Endpoint(), "order_srv.rpc.Find"
	}, opts)

	c.get = kitx.GRPCClientEndpoint(instancer, func(conn *grpc.ClientConn) (endpoint.Endpoint, string) {
		return grpctransport.NewClient(
			conn,
			"pb.UserSvc",
			"Get",
			encodeGetRequest,
			decodeGetResponse,
			orderpb.GetReply{},
			options...,
		).Endpoint(), "order_srv.rpc.Get"
	}, opts)

	return c
}
