package api

import (
	"errors"

	"connectrpc.com/connect"
	"example.com/internal/authn"
	"example.com/internal/service"
)

func AsConnectError(err error) *connect.Error {
	if err == nil {
		return nil
	}

	// https://connectrpc.com/docs/protocol#error-codes

	if errors.Is(err, service.ErrUnauthenticated) {
		return connect.NewError(connect.CodeUnauthenticated, err)
	}

	if errors.Is(err, authn.ErrUnauthenticated) {
		return connect.NewError(connect.CodeUnauthenticated, err)
	}

	if errors.Is(err, service.ErrInsufficientBalance) {
		return connect.NewError(connect.CodeFailedPrecondition, err)
	}

	return connect.NewError(connect.CodeInternal, err)
}
