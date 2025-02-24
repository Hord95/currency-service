package main

import (
	currencyClient "github.com/vctrl/currency-service/currency/internal/clients/currency"
	"github.com/vctrl/currency-service/currency/internal/config"
	"github.com/vctrl/currency-service/currency/internal/db"
	"github.com/vctrl/currency-service/currency/internal/handler"
	"github.com/vctrl/currency-service/currency/internal/repository"
	"github.com/vctrl/currency-service/currency/internal/service"
	"github.com/vctrl/currency-service/pkg/currency"

	"flag"
	"fmt"
	"log"
	"net"

	"go.uber.org/zap"

	"google.golang.org/grpc"
)

// TODO:
// - Добавить run() error по аналогии с migrator
// - Вместо логов - возвращать ошибки

func main() {
	configPath := flag.String("config", "./config", "path to the config file")

	flag.Parse()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("init logger: %v", err)
	}

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	db, _, err := db.NewDatabaseConnection(cfg.Database)
	if err != nil {
		log.Fatalf("error init database connection: %v", err)
	}

	repo, err := repository.NewCurrency(db)
	if err != nil {
		log.Fatalf("error init exchange rate repository: %v", err)
	}

	client, err := currencyClient.New(cfg.API, logger)
	if err != nil {
		log.Fatalf("error creating currency client: %v", err)
	}

	svc := service.NewCurrency(repo, client, logger)

	currencyServer := handler.NewCurrencyServer(svc, logger)

	if err := startGRPCServer(cfg, currencyServer); err != nil {
		log.Fatalf("Error starting GRPC server: %s", err)
	}
}

func startGRPCServer(cfg config.AppConfig, srv handler.CurrencyServer) error {
	lis, err := net.Listen("tcp", ":"+cfg.Service.ServerPort)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	s := grpc.NewServer()
	currency.RegisterCurrencyServiceServer(s, srv)

	log.Printf("gRPC server is listening on :%s", cfg.Service.ServerPort)
	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}
