package grpcapi

import (
	"github.com/cross-team/clublink-backend/app/adapter/grpcapi/proto"
	"github.com/short-d/app/fw/rpc"
	"google.golang.org/grpc"
)

var _ rpc.API = (*Short)(nil)

// Short provides an efficient way for remote systems to interact with Short backend.
type Short struct {
	metaTagServer proto.MetaTagServiceServer
}

// RegisterServers registers gRPC servers that handle user requests.
func (s Short) RegisterServers(server *grpc.Server) {
	proto.RegisterMetaTagServiceServer(server, s.metaTagServer)
}

// NewShort creates Short.
func NewShort(metaTagServer proto.MetaTagServiceServer) Short {
	return Short{metaTagServer: metaTagServer}
}
