package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

type Application struct {
	server  *grpc.Server
	log     *zerolog.Logger
	address string
}

func NewApplication(
	server *grpc.Server,
	log *zerolog.Logger,
	address string,
) Application {
	return Application{
		server:  server,
		log:     log,
		address: address,
	}
}

func (app Application) RunServer() {
	listener, err := net.Listen("tcp", app.address)
	if err != nil {
		app.log.Fatal().Err(err).Msg("failed to create listener")
	}

	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-sigint

		app.server.GracefulStop()

		close(idleConnsClosed)
	}()

	app.log.Info().Str("address", app.address).Msg("gRPC server is running")

	if err := app.server.Serve(listener); err != nil {
		app.log.Fatal().Err(err).Msg("gRPC server error")
	}

	<-idleConnsClosed
	app.log.Info().Msg("gRPC server shutdown gracefully")
}
