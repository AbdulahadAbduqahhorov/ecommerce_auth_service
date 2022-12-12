package main

import (
	"fmt"
	"net"

	"github.com/AbdulahadAbduqahhorov/gRPC/Ecommerce/ecommerce_auth_service/config"
	"github.com/AbdulahadAbduqahhorov/gRPC/Ecommerce/ecommerce_auth_service/genproto/auth_service"
	"github.com/AbdulahadAbduqahhorov/gRPC/Ecommerce/ecommerce_auth_service/service"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/gommon/log"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

func main() {
	var err error
	cfg := config.Load()
	connP := fmt.Sprintf(
		"host=%v port=%v user=%v password=%v dbname=%v sslmode=disable",
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDatabase)

	db, err := sqlx.Connect("postgres", connP)
	if err != nil {
		panic(err)
	}
	authService := service.NewAuthService(cfg, db)

	lis, err := net.Listen("tcp", cfg.GrpcPort)
	if err != nil {
		log.Error("error while listening: %v", err)
		return
	}
	service := grpc.NewServer()
	auth_service.RegisterAuthServiceServer(service, authService)
	if err := service.Serve(lis); err != nil {
		log.Error("error while listening: %v", err)
	}

}
