package grpc

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
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

func NewClient(instancer sd.Instancer, opts ...kitx.Option) *Client {
	c := &Client{}

	c.find = kitx.GRPCClientEndpoint(instancer, func(conn *grpc.ClientConn) (endpoint.Endpoint, string) {
		return grpctransport.NewClient(
			conn,
			"pb.OrderSvc",
			"Find",
			encodeFindRequest,
			decodeFindResponse,
			orderpb.FindReply{},
			//append(options, grpctransport.ClientBefore(opentracing.ContextToGRPC(otTracer, logger)))...,
		).Endpoint(), "order_srv.rpc.Find"
	}, opts...)

	// c.get = kitx.GRPCClientEndpoint(instancer, func(conn *grpc.ClientConn) (endpoint.Endpoint, string) {
	// 	return grpctransport.NewClient(
	// 		conn,
	// 		"pb.UserSvc",
	// 		"Get",
	// 		encodeGetRequest,
	// 		decodeGetResponse,
	// 		userpb.GetReply{},
	// 		//append(options, grpctransport.ClientBefore(opentracing.ContextToGRPC(otTracer, logger)))...,
	// 	).Endpoint(), "user_srv.rpc.Get"
	// }, opts...)

	return c
}
