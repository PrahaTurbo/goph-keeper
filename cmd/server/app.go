package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

// Application holds the gRPC server, logger and the address of the server.
type Application struct {
	server  *grpc.Server
	log     *zerolog.Logger
	address string
}

// NewApplication is a constructor function for Application.
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

// RunServer is a method that starts the server on the provided address.
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
