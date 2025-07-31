package handler

import (
	"context"

	base "github.com/SangBejoo/Template/gen/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"go.elastic.co/apm/v2"
)
func (h *baseHandler) HealthCheck(ctx context.Context, request *emptypb.Empty) (response *base.MessageStatusResponse, err error) {

	span, ctx := apm.StartSpan(ctx, "transport.HealthCheck", "transport.internal")
	defer span.End()

	response = &base.MessageStatusResponse{
		Status:  "OK",
		Message: "Service is healthy",
	}

	if err != nil {
		apm.CaptureError(ctx, err).Send()
	}
	return response, err
}
