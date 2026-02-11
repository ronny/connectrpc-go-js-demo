package api

import (
	"context"

	"connectrpc.com/connect"
	demov1 "example.com/gen/myorg/demo/v1"
	"example.com/internal/authn"
	"example.com/internal/service"
)

func (h *Handler) CreateTransfer(
	ctx context.Context,
	req *connect.Request[demov1.CreateTransferRequest],
) (*connect.Response[demov1.CreateTransferResponse], error) {
	user, err := authn.RequireUser(ctx)
	if err != nil {
		return nil, AsConnectError(err)
	}

	_, err = h.svc.CreateTransfer(ctx, service.CreateTransferRequest{
		SenderEmail:    user.Email,
		RecipientEmail: service.Email(req.Msg.GetRecipientEmail()),
		AmountKoinu:    service.Koinu(req.Msg.GetAmountKoinu()),
	})
	if err != nil {
		return nil, AsConnectError(err)
	}

	resp := new(demov1.CreateTransferResponse)

	return connect.NewResponse(resp), nil
}
