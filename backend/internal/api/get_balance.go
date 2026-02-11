package api

import (
	"context"

	"connectrpc.com/connect"
	demov1 "example.com/gen/myorg/demo/v1"
	"example.com/internal/authn"
	"example.com/internal/service"
)

func (h *Handler) GetBalance(
	ctx context.Context,
	req *connect.Request[demov1.GetBalanceRequest],
) (*connect.Response[demov1.GetBalanceResponse], error) {
	user, err := authn.RequireUser(ctx)
	if err != nil {
		return nil, AsConnectError(err)
	}

	svcResp, err := h.svc.GetBalance(ctx, service.GetBalanceRequest{
		Email: user.Email,
	})
	if err != nil {
		return nil, AsConnectError(err)
	}

	resp := new(demov1.GetBalanceResponse)
	resp.SetKoinu(int64(svcResp.Koinu))

	return connect.NewResponse(resp), nil
}
