package main

import (
	"net"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/zero-shubham/authsvc/internal"
	authrpc "github.com/zero-shubham/authsvc/transport/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	listener, err := net.Listen("tcp", "0.0.0.0:50050")
	if err != nil {
		log.Fatal().Err(err).Msg("error while listening on - :50050")
	}

	log.Info().Msg("listening on - :50050")

	s := grpc.NewServer()
	reflection.Register(s)

	authsvc := &internal.Server{}
	authrpc.RegisterAuthServiceServer(s, authsvc)

	if err := s.Serve(listener); err != nil {
		log.Fatal().Err(err).Msg("failed to serve")
	}
}
