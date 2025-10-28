package main

import (
	"errors"
	"fmt"
	"io"

	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

func newConnectionsContainer(
	config *config,
	_ *log.Logger,
	multiCloser *multiCloser,
) (container *connectionsContainer, err error) {
	containerBuilder := func() error {
		container = &connectionsContainer{}

		db, err := initMySQL(config)
		if err != nil {
			return fmt.Errorf("failed to init DB for migrations: %w", err)
		}
		defer db.Close()

		if err = applyMigrations(db.DB, pathToMigrations); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
		log.Infof("Migrations applied successfully")
		multiCloser.Add(db)
		container.db = db

		// TODO: это конекшены к другим сервисам (в данном случае - gRPC)
		testConnection, err := grpc.NewClient(
			config.TestGRPCAddress,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
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
	db             *sqlx.DB
	testConnection grpc.ClientConnInterface
}

func initMySQL(cfg *config) (db *sqlx.DB, err error) {
	db, err = sqlx.Connect("mysql", cfg.buildDSN())
	if err != nil || db == nil {
		return nil, fmt.Errorf("failed to connect to MySQL: %w", err)
	}
	db.SetMaxOpenConns(cfg.DBMaxConn)

	defer func() {
		if err != nil {
			db.Close()
		}
	}()

	return db, nil
}
