package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/vkotsiuba99/NeoHome/back/internal/service/alert"
	"github.com/vkotsiuba99/NeoHome/back/internal/service/auth"
	"github.com/vkotsiuba99/NeoHome/back/internal/service/device"
	"github.com/vkotsiuba99/NeoHome/back/internal/service/telemetry"
	"github.com/vkotsiuba99/NeoHome/back/internal/storage"
	"github.com/vkotsiuba99/NeoHome/back/internal/storage/cassandra"
	"github.com/vkotsiuba99/NeoHome/back/internal/transport/http/handler"
	"github.com/vkotsiuba99/NeoHome/back/internal/transport/mqtt"
	"github.com/vkotsiuba99/NeoHome/back/pkg/logger"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := Load(ctx)
	if err != nil {
		log.Printf("load config failed: %v", err)
		return
	}

	rootLogger, err := logger.New(cfg.Logger)
	if err != nil {
		log.Printf("logger init failed: %v", err)
		return
	}

	appLogger := rootLogger.With("component", "app")
	httpLogger := rootLogger.With("component", "http")
	dbLogger := rootLogger.With("component", "storage")

	database, err := cassandra.New(ctx, cfg.DB, dbLogger)
	if err != nil {
		log.Printf("init repository failed: %v", err)
		return
	}
	defer func() {
		if closeErr := database.Close(); closeErr != nil {
			dbLogger.Warn("failed to close cassandra connection", "error", closeErr.Error())
		}
	}()

	authRepository := storage.NewAuthRepo(database.DB())
	deviceRepository := storage.NewDeviceRepo(database.DB())
	telemetryRepository := storage.NewTelemetryRepo(database.DB())
	alertRepository := storage.NewAlertRepo(database.DB())

	authService := auth.New(authRepository, cfg.Service, appLogger)
	deviceService := device.New(authRepository, deviceRepository, deviceRepository, cfg.Service, appLogger)
	telemetryService := telemetry.New(deviceRepository, deviceRepository, telemetryRepository, alertRepository, cfg.Service, appLogger)
	alertService := alert.New(deviceRepository, alertRepository, appLogger)

	mqttConsumer, err := mqtt.Start(ctx, cfg.MQTT, telemetryService, rootLogger.With("component", "mqtt"))
	if err != nil {
		log.Printf("mqtt consumer init failed: %v", err)
		return
	}
	defer func() {
		if mqttConsumer != nil {
			mqttConsumer.Close()
		}
	}()

	hand := handler.New(authService, deviceService, telemetryService, alertService, cfg.MQTT.TopicTelemetry, httpLogger)
	router := handler.NewRouter(hand, cfg.HTTP, httpLogger, cfg.Service.JWTSecret)

	serverAddress := fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port)
	httpServer := &http.Server{
		Addr:              serverAddress,
		Handler:           router,
		ReadTimeout:       cfg.HTTP.ReadTimeout,
		ReadHeaderTimeout: cfg.HTTP.ReadHeaderTimeout,
		WriteTimeout:      cfg.HTTP.WriteTimeout,
		IdleTimeout:       cfg.HTTP.IdleTimeout,
		MaxHeaderBytes:    cfg.HTTP.MaxHeaderBytes,
	}

	go func() {
		httpLogger.Info("http server starting", "addr", serverAddress)
		if serveErr := httpServer.ListenAndServe(); serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
			httpLogger.Error("http server failed", "error", serveErr.Error())
			stop()
		}
	}()

	<-ctx.Done()
	httpLogger.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
		return
	}

	appLogger.Info("application stopped gracefully")
}
