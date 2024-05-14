package main

import (
	"context"
	"log/slog"
	"net"

	"github.com/DarkReduX/productservice/internal/config"
	"github.com/DarkReduX/productservice/internal/handler"
	"github.com/DarkReduX/productservice/internal/repository"
	"github.com/DarkReduX/productservice/internal/service"
	"github.com/DarkReduX/productservice/protobuf"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		slog.Error("Failed to initialize configuration: ", slog.String("error", err.Error()))
		return
	}

	ctx := context.Background()

	dbPool, err := pgxpool.New(ctx, cfg.URL)
	if err != nil {
		slog.Error("Failed to initialize database: ", slog.String("error", err.Error()))
		return
	}

	productRep := repository.NewProduct(dbPool)

	productSvc := service.NewProduct(productRep)

	productHandler := handler.NewProduct(productSvc)

	lis, err := net.Listen("tcp", cfg.ListenAddress)
	if err != nil {
		slog.Error("Failed to start listen: ", slog.String("error", err.Error()))
		return
	}

	grpcServer := grpc.NewServer()

	protobuf.RegisterProductServiceServer(grpcServer, productHandler)

	if err = grpcServer.Serve(lis); err != nil {
		slog.Error("Failed gRPC server: ", slog.String("error", err.Error()))
		return
	}
}
