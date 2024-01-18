package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ponyjackal/go-microservice-boilerplate/internal/adapters/database"
	"github.com/ponyjackal/go-microservice-boilerplate/internal/adapters/database/migrations"
	"github.com/ponyjackal/go-microservice-boilerplate/internal/adapters/database/seeds"
	"github.com/ponyjackal/go-microservice-boilerplate/internal/app/routers"
	"github.com/ponyjackal/go-microservice-boilerplate/internal/domain/services"
	server "github.com/ponyjackal/go-microservice-boilerplate/internal/grpc"
	"github.com/ponyjackal/go-microservice-boilerplate/pkg/config"
	"github.com/ponyjackal/go-microservice-boilerplate/pkg/logger"
	// product
)

// @BasePath /api/v1
func main() {
	// init timezone and db
	initDB()
	defer cleanUp()
	// Set up shutdownCh and wg
	shutdownCh := make(chan struct{})
	var wg sync.WaitGroup

	/* service */
	tagService := services.NewTagService()

	// setup router
	router := routers.SetupRoute(tagService)

	// Start the Gin server concurrently in a Goroutine
	serverErrCh := make(chan error)
	go func() {
		serverErrCh <- router.Run(config.ServerConfig())
	}()
	// start grpc server
	server.StartServer(tagService)

	// Graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	// Wait for a signal or a server error
	select {
	case sig := <-sigCh:
		logger.Infof("Received signal: %v. Shutting down...", sig)
	case err := <-serverErrCh:
		logger.Infof("Server error: %v. Shutting down...", err)
	}

	close(shutdownCh) // signal all goroutines to stop
	wg.Wait()         // wait for all goroutines to stop
}

func initDB() {
	//set timezone
	os.Setenv("SERVER_TIMEZONE", "Asia/Tokyo")
	loc, _ := time.LoadLocation(os.Getenv("SERVER_TIMEZONE"))
	time.Local = loc
	// setup db
	if err := config.SetupConfig(); err != nil {
		logger.Fatalf("config SetupConfig() error: %s", err)
	}
	masterDSN, replicaDSN := config.DbConfiguration()

	if err := database.DbConnection(masterDSN, replicaDSN); err != nil {
		logger.Fatalf("database DbConnection error: %s", err)
	}
	// run db migration
	migrations.Migrate()
	// Check if seed data exists
	if !seeds.IsSeedDataExists() {
		// Run db seed
		seeds.SeedData()
	}
}

func cleanUp() {
}
