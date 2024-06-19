package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/lidofinance/finding-forwarder/internal/app/worker"

	"github.com/go-chi/chi/v5"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sync/errgroup"

	nc "github.com/lidofinance/finding-forwarder/internal/connectors/nats"

	"github.com/lidofinance/finding-forwarder/internal/app/server"
	"github.com/lidofinance/finding-forwarder/internal/connectors/logger"
	"github.com/lidofinance/finding-forwarder/internal/connectors/metrics"
	"github.com/lidofinance/finding-forwarder/internal/env"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	defer stop()
	g, gCtx := errgroup.WithContext(ctx)

	cfg, envErr := env.Read("")
	if envErr != nil {
		fmt.Println("Read env error:", envErr.Error())
		return
	}

	log := logger.New(&cfg.AppConfig)

	natsClient, natsErr := nc.New(&cfg.AppConfig, log)
	if natsErr != nil {
		fmt.Println("Could not connect to nats error:", natsErr.Error())
		return
	}
	defer natsClient.Close()

	js, jetStreamErr := jetstream.New(natsClient)
	if jetStreamErr != nil {
		fmt.Println("Could not connect to jetStream error:", jetStreamErr.Error())
		return
	}

	s, createStreamErr := js.CreateStream(gCtx, jetstream.StreamConfig{
		Name:     cfg.AppConfig.NatsStreamName,
		Subjects: []string{fmt.Sprintf(`%s.*`, cfg.AppConfig.NatsStreamName)},
	})

	if createStreamErr != nil && !errors.Is(createStreamErr, nats.ErrStreamNameAlreadyInUse) {
		fmt.Println("Could not create FINDINGS stream error:", createStreamErr.Error())
		return
	}

	log.Info(fmt.Sprintf(`started %s worker`, cfg.AppConfig.Name))

	r := chi.NewRouter()
	promRegistry := prometheus.NewRegistry()
	metricsStore := metrics.New(promRegistry, cfg.AppConfig.MetricsPrefix, cfg.AppConfig.Name, cfg.AppConfig.Env)

	services := server.NewServices(&cfg.AppConfig, metricsStore)
	app := server.New(&cfg.AppConfig, log, metricsStore, js, natsClient)

	app.Metrics.BuildInfo.Inc()
	app.RegisterWorkerRoutes(r)

	alertWorker := worker.NewWorker(log, metricsStore, s, services.Telegram, services.OpsGenia, services.Discord)
	if wrkErr := alertWorker.Run(gCtx, g); wrkErr != nil {
		fmt.Println("Could not start alertWorker error:", wrkErr.Error())
		return
	}

	app.RunHTTPServer(gCtx, g, cfg.AppConfig.Port, r)

	if err := g.Wait(); err != nil {
		log.Error(err.Error())
	}

	fmt.Println(`Main done`)
}