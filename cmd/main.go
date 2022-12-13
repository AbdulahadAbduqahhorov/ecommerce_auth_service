package main

import (
	"fmt"
	"net"

	"github.com/AbdulahadAbduqahhorov/gRPC/Ecommerce/ecommerce_auth_service/config"
	"github.com/AbdulahadAbduqahhorov/gRPC/Ecommerce/ecommerce_auth_service/genproto/auth_service"
	"github.com/AbdulahadAbduqahhorov/gRPC/Ecommerce/ecommerce_auth_service/pkg/logger"
	"github.com/AbdulahadAbduqahhorov/gRPC/Ecommerce/ecommerce_auth_service/service"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

func main() {
	var err error
	cfg := config.Load()
	loggerLevel := logger.LevelDebug
	log := logger.NewLogger("auth_service", loggerLevel)
	defer logger.Cleanup(log)
	connP := fmt.Sprintf(
		"host=%v port=%v user=%v password=%v dbname=%v sslmode=disable",
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDatabase)

	db, err := sqlx.Connect("postgres", connP)
	if err != nil {
		log.Panic("postgres.NewPostgres", logger.Error(err))
	}
	authService := service.NewAuthService(log,cfg, db)

	lis, err := net.Listen("tcp", cfg.GrpcPort)
	if err != nil {
		log.Panic("net.Listen", logger.Error(err))
	}
	service := grpc.NewServer()
	auth_service.RegisterAuthServiceServer(service, authService)

	log.Info("GRPC: Server being started...", logger.String("port", cfg.GrpcPort))

	if err := service.Serve(lis); err != nil {
		log.Panic("grpcServer.Serve", logger.Error(err))
	}

}
