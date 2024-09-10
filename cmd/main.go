package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725721384-team-77753/zadanie-6105/internal/config"
	httphandler "git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725721384-team-77753/zadanie-6105/internal/handlers/http"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725721384-team-77753/zadanie-6105/internal/repository/postgres"
	"git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725721384-team-77753/zadanie-6105/internal/server"
	service "git.codenrock.com/avito-testirovanie-na-backend-1270/cnrprod1725721384-team-77753/zadanie-6105/internal/services"
)

func main() {
	log := server.SetupLogger()

	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	// Migrate: add new tables if not exists
	if err := server.Migrate(&cfg.PG, log); err != nil {
		log.Fatal(err)
	}

	// Layers
	repo, err := postgres.New(&cfg.PG)
	if err != nil {
		log.Fatalf("unable to connect to pool: %v", err)
	}

	srv := service.New(repo, log)
	handler := httphandler.New(srv, log)

	// Server
	app := server.New(handler.InitRoutes(), &cfg.Server)
	go func() {
		log.Infof("the server is starting on %v", cfg.Server.Address)

		if err := app.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Errorf("error occured while running http server: %s", err.Error())
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGTERM, syscall.SIGINT)
	<-done

	if err := app.Shutdown(context.Background()); err != nil {
		log.Errorf("error occured on server shutting down: %s", err)
	}

	log.Info("the server is shut down")
}
