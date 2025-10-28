package main

import (
	"fmt"
	"time"

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

	DBHost     string `envconfig:"db_host" default:"localhost"`
	DBPort     string `envconfig:"db_port"`
	DBName     string `envconfig:"db_name"`
	DBUser     string `envconfig:"db_user"`
	DBPassword string `envconfig:"db_password"`
	DBMaxConn  int    `envconfig:"db_max_conn"`

	TestGRPCAddress string `envconfig:"test_grpc_address" default:"test:8081"`
}

func (c *config) buildDSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=%s",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName,
		time.UTC.String(),
	)
}
