package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/edaywalid/reverse-proxy/internal/app"
	"github.com/edaywalid/reverse-proxy/pkg/utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Info().Msg("Starting server...")

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)

	app, err := app.New(utils.CONFIG_FILE)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create app")
	}

	go func() {
		log.Info().Msgf("Starting HTTP redirect server on port %d", app.Config.HTTPPort)
		port := fmt.Sprintf(":%d", app.Config.HTTPPort)
		if err := http.ListenAndServe(port, http.HandlerFunc(app.Handlers.ProxyHandler.HttpRedirect)); err != nil {
			log.Fatal().Err(err).Msg("Failed to start HTTP redirect server")
		}
	}()

	http.HandleFunc("/", app.Handlers.ProxyHandler.Handler)
	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", app.Config.HTTPSPort),
	}

	go func() {
		log.Info().Msgf("Starting server on port %d", app.Config.HTTPSPort)
		if err := srv.ListenAndServeTLS(app.Config.CertFile, app.Config.KeyFile); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	<-stopChan
	log.Info().Msg("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Failed to shutdown server")
	}
	log.Info().Msg("Server shutdown")
}
