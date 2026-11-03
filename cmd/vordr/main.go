package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/docker/docker/client"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sync/errgroup"

	"github.com/st0o0/vordr/internal/config"
	"github.com/st0o0/vordr/internal/metrics"
	"github.com/st0o0/vordr/internal/scanner"
	"github.com/st0o0/vordr/internal/scheduler"
	"github.com/st0o0/vordr/internal/store"
	"github.com/st0o0/vordr/internal/trivy"
	"github.com/st0o0/vordr/internal/watcher"
	"github.com/st0o0/vordr/internal/worker"
)

var version = "dev"

func main() {
	if len(os.Args) > 1 && os.Args[1] == "healthcheck" {
		os.Exit(runHealthcheck())
	}

	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Println(version)
		os.Exit(0)
	}

	cfg, err := config.Load(os.Getenv)
	if err != nil {
		fmt.Fprintf(os.Stderr, "configuration error: %v\n", err)
		os.Exit(1)
	}

	logger := newLogger(cfg)

	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logger.Error("failed to create docker client", "error", err)
		os.Exit(1)
	}
	defer dockerClient.Close()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	startTime := time.Now()
	scanCh := make(chan scanner.ImageRef, 2*cfg.ScanWorkers)

	st := store.New()
	trivyClient := trivy.NewCLI(cfg.TrivyCacheDir, cfg.TrivyTimeout)
	pool := worker.NewPool(cfg.ScanWorkers, scanCh, st, trivyClient, logger)
	w := watcher.New(dockerClient, st, scanCh, cfg.ScanStagger, logger)
	sched := scheduler.New(cfg.ScanInterval, st, scanCh, logger)

	registry := prometheus.NewRegistry()
	collector := metrics.NewCollector(st, &pool.ActiveWorkers, func() int { return len(scanCh) })
	registry.MustRegister(collector)

	srv := metrics.NewServer(cfg.ListenAddr, registry, version, startTime, logger)

	logger.Info("vordr starting",
		"version", version,
		"trivy_cache_dir", cfg.TrivyCacheDir,
		"scan_interval", cfg.ScanInterval,
		"scan_workers", cfg.ScanWorkers,
		"metrics_addr", cfg.ListenAddr,
	)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		srv.Start(ctx)
		<-ctx.Done()
		return nil
	})

	g.Go(func() error { return w.Run(ctx) })
	g.Go(func() error { return sched.Run(ctx) })
	g.Go(func() error { return pool.Run(ctx) })

	if err := g.Wait(); err != nil {
		logger.Error("fatal error", "error", err)
		os.Exit(1)
	}

	logger.Info("shutdown complete")
}

func newLogger(cfg *config.Config) *slog.Logger {
	opts := &slog.HandlerOptions{Level: cfg.LogLevel}
	var handler slog.Handler
	if cfg.LogFormat == "json" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}
	return slog.New(handler)
}
