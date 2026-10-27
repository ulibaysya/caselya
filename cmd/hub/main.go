package main

import (
	"context"
	"log/slog"
	"net"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	phuslog "github.com/phuslu/log"
	itemtagh "github.com/ulibaysya/caselya/internal/modules/itemtag/handler"
	itemtagsrvc "github.com/ulibaysya/caselya/internal/modules/itemtag/service"
	itemtagstrg "github.com/ulibaysya/caselya/internal/modules/itemtag/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	logger := slog.New(phuslog.SlogNewJSONHandler(os.Stderr, &slog.HandlerOptions{
		// Level: slog.LevelDebug,
		Level: slog.LevelInfo,
	}))

	db, err := pgxpool.New(context.Background(), "dbname=caselya")
	if err != nil {
		phuslog.Fatal().Err(err).Msg("")
	}

	serv := grpc.NewServer()

	itemtagStorage := itemtagstrg.New(logger, db)
	itemtagService := itemtagsrvc.New(logger, itemtagStorage)
	itemtagh.Register(serv, itemtagService)

	reflection.Register(serv)

	l, err := net.Listen("tcp", ":50051")
	if err != nil {
		phuslog.Fatal().Err(err).Msg("")
	}
	defer l.Close()
	phuslog.Info().Msg("listening")
	phuslog.Fatal().Err(serv.Serve(l)).Msg("")
}
