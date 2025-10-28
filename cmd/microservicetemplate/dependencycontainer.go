package main

import "github.com/jmoiron/sqlx"

// TODO: добавить зависимости

func newDependencyContainer(
	_ *config,
	connContainer *connectionsContainer,
) (*dependencyContainer, error) {
	return &dependencyContainer{
		db: connContainer.db,
	}, nil
}

type dependencyContainer struct {
	db *sqlx.DB
}
