package main

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

func parseEnv() (*config, error) {
	c := new(config)
	if err := envconfig.Process(appID, c); err != nil {
		return nil, errors.Wrap(err, "failed to parse env")
	}
	return c, nil
}

type config struct {
	LogLevel string `envconfig:"log_level" default:"info"`

	ServeGRPCAddress string `envconfig:"serve_grpc_address" default:":8081"`

	TestGRPCAddress string `envconfig:"test_grpc_address" default:"test:8081"`
}
