package main

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"

	"microservicetemplate"
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
			return err
		}
		multiCloser.Add(db)
		container.db = db

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
	db             *sqlx.DB
	testConnection grpc.ClientConnInterface
}

func initMySQL(cfg *config) (db *sqlx.DB, err error) {
	db, err = sqlx.Connect("mysql", cfg.buildDSN())
	if err != nil || db == nil {
		return nil, fmt.Errorf("failed to connect to MySQL after retries: %w", err)
	}

	db.SetMaxOpenConns(cfg.DBMaxConn)

	go func() {
		if err != nil {
			db.Close()
		}
	}()

	if err := applyMigrations(db.DB); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	return db, nil
}

func applyMigrations(db *sql.DB) error {
	d, err := iofs.New(microservicetemplate.Migrations, "data/mysql/migrations")
	if err != nil {
		return fmt.Errorf("failed to create migration source: %w", err)
	}

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return fmt.Errorf("failed to create MySQL driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", d, "mysql", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate up failed: %w", err)
	}

	return nil
}
