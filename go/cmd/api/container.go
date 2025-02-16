package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/ahmedalhulaibi/loggy"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.opencensus.io/plugin/ocgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	helloworldv1 "github.com/ahmedalhulaibi/hello-world-api/internal/gen/helloworld/v1"
	"github.com/ahmedalhulaibi/hello-world-api/internal/greeter"
	"github.com/ahmedalhulaibi/hello-world-api/internal/grpcutil/interceptors/instanceid"
	logmw "github.com/ahmedalhulaibi/hello-world-api/internal/grpcutil/interceptors/log"
	"github.com/ahmedalhulaibi/hello-world-api/internal/grpcutil/interceptors/requestid"
	"github.com/ahmedalhulaibi/hello-world-api/internal/grpcutil/interceptors/userid"
	httputilgrpcgateway "github.com/ahmedalhulaibi/hello-world-api/internal/httputil/grpcgateway"
	"github.com/ahmedalhulaibi/hello-world-api/internal/tracing"
)

type container struct {
	config *Config

	state struct {
		logger *loggy.Logger

		grpcServer    *grpc.Server
		gatewayRouter *runtime.ServeMux
		gatewayServer *http.Server

		grpcListener    net.Listener
		gatewayListener net.Listener

		greeterService helloworldv1.GreeterServiceServer
	}

	once struct {
		logger, grpcServer, gatewayRouter, gatewayServer, grpcListener, gatewayListener, greeterService sync.Once
	}
}

func newContainer(config *Config) *container {
	return &container{
		config: config,
	}
}

func (c *container) greeterService() helloworldv1.GreeterServiceServer {
	c.once.greeterService.Do(func() {
		c.state.greeterService = greeter.NewGreeter("Hello, %s! Ya filthy animal.")
	})

	return c.state.greeterService
}

func (c *container) logger() *loggy.Logger {
	c.once.logger.Do(func() {
		var lc zap.Config
		switch {
		case c.config.Dev:
			lc = zap.NewDevelopmentConfig()
		default:
			lc = zap.NewProductionConfig()
		}

		logger, err := lc.Build(zap.AddCallerSkip(1))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		serverLogger := loggy.New(logger.Sugar())

		c.state.logger = &serverLogger
	})

	return c.state.logger
}

func (c *container) grpcServer() *grpc.Server {
	c.once.grpcServer.Do(func() {
		c.state.grpcServer = grpc.NewServer(
			grpc.StatsHandler(&ocgrpc.ServerHandler{}),
			grpc.ChainUnaryInterceptor(
				requestid.RequestIdUnaryServerInterceptor(c.logger()),
				instanceid.InstanceIdUnaryServerInterceptor(c.logger(), c.config.Server.InstanceID),
				userid.UserIdUnaryServerInterceptor(c.logger()),
				tracing.NewOpenCensusTraceInterceptor(c.logger()),
				logmw.LoggerUnaryServerInterceptor(c.logger()),
			),
		)

		helloworldv1.RegisterGreeterServiceServer(c.state.grpcServer, c.greeterService())
		reflection.Register(c.state.grpcServer)
	})

	return c.state.grpcServer
}

func (c *container) gatewayRouter() *runtime.ServeMux {
	c.once.gatewayRouter.Do(func() {
		c.state.gatewayRouter = runtime.NewServeMux(
			runtime.WithIncomingHeaderMatcher(httputilgrpcgateway.CustomMatcher),
		)

		ctx := context.Background()
		conn, err := grpc.NewClient(
			fmt.Sprintf("0.0.0.0%s", c.config.Server.GRPCAddr),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			c.logger().Fatalw(ctx, "gateway-router", "err", err)
		}

		err = helloworldv1.RegisterGreeterServiceHandler(
			ctx,
			c.state.gatewayRouter,
			conn,
		)
		if err != nil {
			c.logger().Fatalw(context.Background(), "gateway-router", "err", err)
		}
	})

	return c.state.gatewayRouter
}

func (c *container) gatewayServer() *http.Server {
	c.once.gatewayServer.Do(func() {
		gatewayRouter := c.gatewayRouter()

		c.state.gatewayServer = &http.Server{
			Addr:         c.config.Server.GatewayAddr,
			ReadTimeout:  c.config.Server.Timeout,
			WriteTimeout: c.config.Server.Timeout,
			Handler:      gatewayRouter,
			// Handler: &ochttp.Handler{
			// 	Handler:     gatewayRouter,
			// 	Propagation: &b3.HTTPFormat{},
			// },
		}
	})

	return c.state.gatewayServer
}

func (c *container) gatewayListener() net.Listener {
	c.once.gatewayListener.Do(func() {
		listener, err := net.Listen("tcp", c.config.Server.GatewayAddr)
		if err != nil {
			c.logger().Fatalw(context.Background(), "gateway-listener", "addr", c.config.Server.GatewayAddr, "err", err)
		}

		c.state.gatewayListener = listener
	})

	return c.state.gatewayListener
}

func (c *container) grpcListener() net.Listener {
	c.once.grpcListener.Do(func() {
		listener, err := net.Listen("tcp", c.config.Server.GRPCAddr)
		if err != nil {
			c.logger().Fatalw(context.Background(), "grpc-listener", "addr", c.config.Server.GRPCAddr, "err", err)
		}

		c.state.grpcListener = listener
	})

	return c.state.grpcListener
}
