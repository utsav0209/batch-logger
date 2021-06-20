package utils

import (
	"context"
	"github.com/go-chi/chi/v5/middleware"
)

type PayloadRequestId struct {
	id string
}

func GetPayloadRequestId(ctx context.Context) PayloadRequestId {
	return PayloadRequestId{middleware.GetReqID(ctx) + "-payload"}
}
