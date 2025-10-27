package main

// TODO: добавить зависимости

func newDependencyContainer(
	_ *config,
	_ *connectionsContainer,
) (*dependencyContainer, error) {
	return &dependencyContainer{}, nil
}

type dependencyContainer struct {
}
