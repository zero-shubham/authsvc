package main

import (
	"context"
	"net"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	pgxUUID "github.com/vgarvardt/pgx-google-uuid/v5"
	"github.com/zero-shubham/authsvc/internal"
	authrpc "github.com/zero-shubham/authsvc/transport/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	MaxIdelConn = 5
	DbUrlEnv    = "DATABASE_URL"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	ctx := context.Background()
	DB_URL := os.Getenv(DbUrlEnv)

	listener, err := net.Listen("tcp", "0.0.0.0:50050")
	if err != nil {
		log.Fatal().Err(err).Msg("error while listening on - :50050")
	}

	log.Info().Msg("listening on - :50050")

	s := grpc.NewServer()
	reflection.Register(s)

	dbConfig, err := pgxpool.ParseConfig(DB_URL)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create a config")
	}
	dbConfig.MaxConns = MaxIdelConn
	dbConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxUUID.Register(conn.TypeMap())
		return nil
	}

	dbConn, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create a config")
	}

	authsvc := internal.NewServer(dbConn)
	authrpc.RegisterAuthServiceServer(s, authsvc)

	if err := s.Serve(listener); err != nil {
		log.Fatal().Err(err).Msg("failed to serve")
	}
}
