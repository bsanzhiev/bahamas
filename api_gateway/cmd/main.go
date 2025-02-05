package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	gateway_services "github.com/bsanzhiev/bahamas/api_gateway/internal/application/services"
	"github.com/bsanzhiev/bahamas/api_gateway/internal/infrastructure/config"
	"github.com/bsanzhiev/bahamas/api_gateway/internal/infrastructure/http/server"
	"google.golang.org/grpc"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	// Init logger
	logger := log.New(os.Stdout, "API Gateway: ", log.LstdFlags|log.Lshortfile)

	// Context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Init gRPC clients
	customerConn, err := grpc.NewClient(cfg.CustomerServiceURL)
	if err != nil {
		logger.Fatalf("Failed to create customers service client: %v", err)
	}
	defer customerConn.Close()

	transactionConn, err := grpc.NewClient(cfg.TransactionServiceURL)
	if err != nil {
		logger.Fatalf("Failed to create transactions service client: %v", err)
	}
	defer transactionConn.Close()

	// Init services
	customerService := gateway_services.NewCustomerService(customerConn)
	transactionService := gateway_services.NewTransactionService(transactionConn)

	// Init HTTP handlers
	customerHandler := http_handler.NewCustomerHandler(customerService, logger)
	transactionHandler := http_handler.NewTransactionHandler(transactionService, logger)

	// Setup middleware
	chain := middleware.Chain(
		middleware.Logger(logger),
		middleware.Authenticate(cfg.AuthSecret),
	)

	// Setup routes
	router := server.NewRouter()
	router.Handle("/api/v1/customers", chain.Then(handlers.HandleCustomers(customerHandler)))
	router.Handle("/api/v1/transactions", chain.Then(handlers.HandleTransactions(transactionHandler)))

	// Setup HTTP server
	srv := server.NewServer(cfg.HTTP.Port, router, logger)

	// Start HTTP server
	go func() {
		logger.Printf("Starting HTTP server on port %s", cfg.HTTP.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("HTTP server failed: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logger.Println("Shutting down server...")

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Printf("HTTP server shutdown failed: %v", err)
	}

	logger.Println("Server exited properly")
}
