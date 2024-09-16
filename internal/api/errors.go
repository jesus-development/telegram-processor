package api

import (
	"errors"
	"google.golang.org/grpc/status"
)

var (
	ErrNotFound           = status.Error(404, "not found")
	ErrSomethingWentWrong = status.Error(500, "something went wrong")
	ErrBadRequest         = status.Error(400, "bad request")
)

var (
	ErrGetTraceId = errors.New("failed to get trace id")
)
