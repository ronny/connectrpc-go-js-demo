package api

import (
	"context"

	"connectrpc.com/connect"
	demov1 "example.com/gen/myorg/demo/v1"
	"example.com/internal/service"
)

func (h *Handler) Login(
	ctx context.Context,
	req *connect.Request[demov1.LoginRequest],
) (*connect.Response[demov1.LoginResponse], error) {
	svcResp, err := h.svc.Login(ctx, service.LoginRequest{
		Email:    service.Email(req.Msg.GetEmail()),
		Password: req.Msg.GetPassword(),
	})
	if err != nil {
		return nil, AsConnectError(err)
	}

	resp := new(demov1.LoginResponse)
	resp.SetAuthToken(string(svcResp.AuthToken))

	return connect.NewResponse(resp), nil
}
