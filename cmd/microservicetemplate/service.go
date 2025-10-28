package main

import (
	"context"
	"net"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"

	api "microservicetemplate/api/server/microservicetemplateinternal"
	"microservicetemplate/pkg/infrastructure/transport"
)

const shutdownTimeout = 30 * time.Second

func service(
	config *config,
	logger *log.Logger,
	closer *multiCloser,
) *cli.Command {
	return &cli.Command{
		Name:  "service",
		Usage: "Runs the gRPC service",
		Action: func(c *cli.Context) error {
			connContainer, err := newConnectionsContainer(config, logger, closer)
			if err != nil {
				return errors.Wrap(err, "failed to init connections")
			}

			container, err := newDependencyContainer(config, connContainer)
			if err != nil {
				return errors.Wrap(err, "failed to init dependencies")
			}
			return startGRPCServer(c.Context, config, logger, container)
		},
	}
}

func startGRPCServer(
	ctx context.Context,
	config *config,
	logger *log.Logger,
	_ *dependencyContainer,
) error {
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(makeGrpcUnaryInterceptor(logger)))

	// TODO: зарегистрировать свой сервер вместо шаблонного
	api.RegisterMicroserviceTemplateInternalServiceServer(grpcServer, transport.NewInternalAPI())

	listener, err := net.Listen("tcp", config.ServeGRPCAddress)
	if err != nil {
		return errors.Wrapf(err, "failed to listen on %s", config.ServeGRPCAddress)
	}
	logger.Infof("gRPC server listening on %s", config.ServeGRPCAddress)

	errCh := make(chan error, 1)
	go func() {
		errCh <- grpcServer.Serve(listener)
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		logger.Infof("Shutdown signal received, stopping gRPC server...")
		shutdownGRPCServer(grpcServer, logger)
		return nil
	}
}

func shutdownGRPCServer(server *grpc.Server, logger *log.Logger) {
	done := make(chan struct{})
	go func() {
		server.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		logger.Infof("gRPC server stopped gracefully")
	case <-time.After(shutdownTimeout):
		logger.Warnf("Graceful shutdown timed out after %v, forcing stop", shutdownTimeout)
		server.Stop()
	}
}

func makeGrpcUnaryInterceptor(logger *log.Logger) grpc.UnaryServerInterceptor {
	loggerInterceptor := transport.MakeLoggerServerInterceptor(logger)
	errorInterceptor := transport.ErrorInterceptor{Logger: logger}
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		resp, err = loggerInterceptor(ctx, req, info, handler)
		return resp, errorInterceptor.TranslateGRPCError(err)
	}
}
