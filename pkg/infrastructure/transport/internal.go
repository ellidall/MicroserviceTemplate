package transport

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	api "microservicetemplate/api/server/microservicetemplateinternal"
)

func NewInternalAPI() api.MicroserviceTemplateInternalServiceServer {
	return &internalAPI{}
}

type internalAPI struct {
}

func (i *internalAPI) Ping(_ context.Context, _ *emptypb.Empty) (*api.PingResponse, error) {
	// TODO implement me
	panic("implement me")
}
