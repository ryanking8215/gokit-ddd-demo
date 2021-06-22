package kitx

import (
	"sync"

	"google.golang.org/grpc"
)

type GRPCConnRef struct {
	conn *grpc.ClientConn
	ref  int32
	mu   sync.Mutex
}

func NewGRPConnRef() *GRPCConnRef {
	return &GRPCConnRef{}
}

func (gcr *GRPCConnRef) Conn(instance string) (*grpc.ClientConn, error) {
	gcr.mu.Lock()
	defer gcr.mu.Unlock()

	if gcr.ref == 0 {
		conn, err := grpc.Dial(instance, grpc.WithInsecure())
		if err != nil {
			return nil, err
		}
		gcr.conn = conn
	}

	gcr.ref++
	return gcr.conn, nil
}

func (gcr *GRPCConnRef) Close() error {
	gcr.mu.Lock()
	defer gcr.mu.Unlock()

	gcr.ref--
	if gcr.ref <= 0 {
		return gcr.conn.Close()
	}
	return nil
}
