package main

import (
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
)

type multiCloser struct {
	closers []io.Closer
}

func (m *multiCloser) Add(c io.Closer) {
	if c != nil {
		m.closers = append(m.closers, c)
	}
}

func (m *multiCloser) Close() error {
	var errs []error
	for _, c := range m.closers {
		if err := c.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) == 0 {
		return nil
	}

	return errors.Join(errs...)
}

func newConnectionsContainer(config *config, multiCloser *multiCloser) (container *connectionsContainer, err error) {
	containerBuilder := func() error {
		container = &connectionsContainer{}

		// TODO: это конекшены к другим сервисам (в данном случае - gRPC)
		testConnection, err := grpc.NewClient(
			config.TestGRPCAddress,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(maxGRPCMsgSize), grpc.MaxCallRecvMsgSize(maxGRPCMsgSize)),
		)
		if err != nil {
			return err
		}

		multiCloser.Add(testConnection)
		container.testConnection = testConnection

		return nil
	}

	return container, containerBuilder()
}

type connectionsContainer struct {
	testConnection grpc.ClientConnInterface
}
