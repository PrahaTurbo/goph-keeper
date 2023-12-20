package main

import (
	"fmt"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/PrahaTurbo/goph-keeper/api/proto"
	"github.com/PrahaTurbo/goph-keeper/internal/server/config"
	"github.com/PrahaTurbo/goph-keeper/internal/server/encryption"
	"github.com/PrahaTurbo/goph-keeper/internal/server/handlers"
	"github.com/PrahaTurbo/goph-keeper/internal/server/interceptors"
	"github.com/PrahaTurbo/goph-keeper/internal/server/jwt"
	"github.com/PrahaTurbo/goph-keeper/internal/server/repository"
	"github.com/PrahaTurbo/goph-keeper/internal/server/repository/pg"
	"github.com/PrahaTurbo/goph-keeper/internal/server/services"
	"github.com/PrahaTurbo/goph-keeper/pkg/logger"
)

func main() {
	log := logger.NewLogger().With().
		Int("pid", os.Getpid()).
		Str("app", "goph-keeper-server").
		Logger()

	cfg := config.LoadConfig()

	pgPool, err := pg.NewPGPool(cfg.PG)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to setup database connection")
	}
	defer pgPool.Close()

	jwtManager := jwt.NewJWTManager(cfg.Server.Secret)
	cryptoSrvc := encryption.NewCryptoService(cfg.Server.Secret)

	authRepo := repository.NewAuthRepository(pgPool)
	secretRepo := repository.NewSecretRepository(pgPool)

	authService := services.NewAuthService(authRepo, &log, jwtManager)
	secretService := services.NewSecretService(secretRepo, &log, cryptoSrvc)

	authHandler := handlers.NewAuthHandler(authService, &log)
	secretHandler := handlers.NewSecretHandler(secretService, &log)

	authInterceptor := interceptors.NewAuthInterceptor(jwtManager)

	creds, err := credentials.NewServerTLSFromFile(cfg.Server.CertPath, cfg.Server.KeyPath)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to make creds for server")
	}

	opts := []grpc.ServerOption{
		grpc.Creds(creds),
		grpc.UnaryInterceptor(authInterceptor.UnaryServerInterceptor),
	}

	server := grpc.NewServer(opts...)

	pb.RegisterAuthServer(server, authHandler)
	pb.RegisterSecretServer(server, secretHandler)

	app := NewApplication(server, &log, fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port))

	app.RunServer()
}
