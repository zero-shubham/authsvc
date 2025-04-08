package internal

import (
	"context"

	authrpc "github.com/zero-shubham/authsvc/transport/grpc"
)

type Server struct {
	authrpc.AuthServiceServer
}

func (s *Server) UserToken(ctx context.Context, in *authrpc.UserTokenRequest) (*authrpc.TokenResponse, error) {
	return &authrpc.TokenResponse{
		AccessToken: "test-token",
	}, nil
}
