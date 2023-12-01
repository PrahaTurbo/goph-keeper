package main

import (
	"fmt"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/PrahaTurbo/goph-keeper/api/proto"
	"github.com/PrahaTurbo/goph-keeper/internal/client/config"
	"github.com/PrahaTurbo/goph-keeper/internal/client/tui"
	"github.com/PrahaTurbo/goph-keeper/pkg/logger"
)

func main() {
	log := logger.NewLogger().With().
		Int("pid", os.Getpid()).
		Str("app", "goph-keeper-client").
		Logger()

	cfg := config.LoadConfig()

	creds, err := credentials.NewClientTLSFromFile(cfg.SSLCertPath, "localhost")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create creds")
	}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", cfg.Host, cfg.Port), grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to setup connection")
	}
	defer conn.Close()

	authClient := pb.NewAuthClient(conn)
	secretsClient := pb.NewSecretClient(conn)

	ui := tui.NewApplication(authClient, secretsClient)

	if err := ui.App.SetRoot(ui.Pages, true).EnableMouse(true).Run(); err != nil {
		log.Fatal().Err(err).Msg("client error")
	}
}
