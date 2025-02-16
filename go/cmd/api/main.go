package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"contrib.go.opencensus.io/exporter/ocagent"
	"go.opencensus.io/trace"
	"golang.org/x/sync/errgroup"
)

func main() {
	config, err := parseConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	container := newContainer(config)

	// TODO: swap out for opentelemetry generic exporter when otel collector is ready
	if oce, err := ocagent.NewExporter(
		ocagent.WithInsecure(),
		ocagent.WithReconnectionPeriod(5*time.Second),
		ocagent.WithAddress(config.Server.OpenCensusAgentHost),
		ocagent.WithServiceName(config.Server.ServiceName),
	); err == nil {
		trace.RegisterExporter(oce)
		trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := run(ctx, container); err != nil {
		container.logger().Errorw(ctx, "run", "err", err)
	}
}

func run(ctx context.Context, c *container) error {
	errg, ctx := errgroup.WithContext(ctx)

	runGRPCServer(ctx, errg, c)
	runGatewayServer(ctx, errg, c)

	return errg.Wait()
}

func runGRPCServer(ctx context.Context, errg *errgroup.Group, c *container) {
	grpcServer := c.grpcServer()

	grpcListener := c.grpcListener()
	grpcAddr := grpcListener.Addr().String()
	c.logger().Infow(ctx, "starting grpc server", "addr", grpcAddr)

	errg.Go(func() error {
		<-ctx.Done()

		sctx, cancel := context.WithTimeout(context.Background(), c.config.Server.ShutdownTimeout)
		defer cancel()

		grpcServer.GracefulStop()
		<-sctx.Done()
		grpcServer.Stop()

		c.logger().Infow(ctx, "grpc server shutdown", "addr", grpcAddr)
		return nil
	})

	errg.Go(func() error {
		if err := grpcServer.Serve(grpcListener); err != nil && err != http.ErrServerClosed {
			return err
		}

		return nil
	})
}

func runGatewayServer(ctx context.Context, errg *errgroup.Group, c *container) {
	gatewayServer := c.gatewayServer()

	gatewayListener := c.gatewayListener()
	gatewayAddr := gatewayListener.Addr().String()
	c.logger().Infow(ctx, "starting gateway server", "addr", gatewayAddr)

	errg.Go(func() error {
		<-ctx.Done()

		sctx, cancel := context.WithTimeout(context.Background(), c.config.Server.ShutdownTimeout)
		defer cancel()

		if err := gatewayServer.Shutdown(sctx); err != nil {
			return fmt.Errorf("gateway shutdown: %w", err)
		}

		c.logger().Infow(ctx, "gateway server shutdown", "addr", gatewayAddr)
		return nil
	})

	errg.Go(func() error {
		if err := gatewayServer.Serve(gatewayListener); err != nil && err != http.ErrServerClosed {
			return err
		}

		return nil
	})
}
