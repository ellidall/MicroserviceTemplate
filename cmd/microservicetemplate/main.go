package main

import (
	"context"
	stdlog "log"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// TODO:  appID используется как префикс для env-переменных

const appID = "microserviceTemplate"

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cnf, err := parseEnv()
	if err != nil {
		stdlog.Fatal(err)
	}
	logger, err := initLogger(cnf.LogLevel)
	if err != nil {
		stdlog.Fatal("failed to initialize logger")
	}

	err = runApp(ctx, cnf, logger)
	switch errors.Cause(err) {
	case nil:
		logger.Infof("call finished")
	default:
		logger.Fatal(err)
	}
}

func runApp(
	ctx context.Context,
	config *config,
	logger *log.Logger,
) (err error) {
	closer := &multiCloser{}
	defer func() {
		if closeErr := closer.Close(); closeErr != nil {
			err = errors.Wrap(err, closeErr.Error())
			if err == nil {
				err = closeErr
			}
		}
	}()

	app := cli.App{
		Name: appID,
		Commands: []*cli.Command{
			service(config, logger, closer),
			migrate(config, logger),
		},
	}

	return app.RunContext(ctx, os.Args)
}

func initLogger(level string) (*log.Logger, error) {
	lvl, err := log.ParseLevel(level)
	if err != nil {
		return nil, err
	}

	logger := log.New()
	logger.SetLevel(lvl)
	logger.SetFormatter(&log.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
	})

	return logger, nil
}
