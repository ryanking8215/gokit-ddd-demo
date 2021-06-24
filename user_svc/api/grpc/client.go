package grpc

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/tracing/zipkin"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"

	"gokit-ddd-demo/lib/kitx"
	userpb "gokit-ddd-demo/user_svc/api/grpc/pb"
	"gokit-ddd-demo/user_svc/svc/user"
)

var _ user.Service = (*Client)(nil)

type Client struct {
	find endpoint.Endpoint
	get  endpoint.Endpoint
}

func (c *Client) Find(ctx context.Context) ([]*user.User, error) {
	rsp, err := c.find(ctx, nil)
	if err != nil {
		return nil, err
	}
	return rsp.([]*user.User), err
}

func (c *Client) Get(ctx context.Context, id int64) (*user.User, error) {
	rsp, err := c.get(ctx, id)
	if err != nil {
		return nil, err
	}
	return rsp.(*user.User), nil
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
			"pb.UserSvc",
			"Find",
			encodeFindRequest,
			decodeFindResponse,
			userpb.FindReply{},
			options...,
		).Endpoint(), "user_srv.rpc.Find"
	}, opts)

	c.get = kitx.GRPCClientEndpoint(instancer, func(conn *grpc.ClientConn) (endpoint.Endpoint, string) {
		return grpctransport.NewClient(
			conn,
			"pb.UserSvc",
			"Get",
			encodeGetRequest,
			decodeGetResponse,
			userpb.GetReply{},
			options...,
		).Endpoint(), "user_srv.rpc.Get"
	}, opts)

	return c
}
