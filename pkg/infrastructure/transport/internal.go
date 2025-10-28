package transport

import (
	"context"

	api "microservicetemplate/api/server/microservicetemplateinternal"
)

func NewInternalAPI() api.MicroserviceTemplateInternalServiceServer {
	return &internalAPI{}
}

type internalAPI struct {
}

func (i *internalAPI) Ping(_ context.Context, _ *api.PingRequest) (*api.PingResponse, error) {
	return &api.PingResponse{
		Message: "pong",
	}, nil
}
