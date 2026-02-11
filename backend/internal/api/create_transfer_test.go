package api_test

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	demov1 "example.com/gen/myorg/demo/v1"
)

func TestCreateTransfer_Success(t *testing.T) {
	client := newTestClient(t)
	voidToken := loginAs(t, client, "void@example.com")
	shibaToken := loginAs(t, client, "shiba@example.com")

	msg := new(demov1.CreateTransferRequest)
	msg.SetRecipientEmail("shiba@example.com")
	msg.SetAmountKoinu(5000)

	_, err := client.CreateTransfer(context.Background(), withAuth(msg, voidToken))
	if err != nil {
		t.Fatal(err)
	}

	// Verify void balance decreased
	voidResp, err := client.GetBalance(context.Background(), withAuth(new(demov1.GetBalanceRequest), voidToken))
	if err != nil {
		t.Fatal(err)
	}

	if voidResp.Msg.GetKoinu() != 1_234_567_890-5000 {
		t.Fatalf("expected void balance %d, got %d", 1_234_567_890-5000, voidResp.Msg.GetKoinu())
	}

	// Verify shiba balance increased
	shibaResp, err := client.GetBalance(context.Background(), withAuth(new(demov1.GetBalanceRequest), shibaToken))
	if err != nil {
		t.Fatal(err)
	}

	if shibaResp.Msg.GetKoinu() != 2_111_222_333+5000 {
		t.Fatalf("expected shiba balance %d, got %d", 2_111_222_333+5000, shibaResp.Msg.GetKoinu())
	}
}

func TestCreateTransfer_MissingAuthToken(t *testing.T) {
	client := newTestClient(t)

	msg := new(demov1.CreateTransferRequest)
	msg.SetRecipientEmail("shiba@example.com")
	msg.SetAmountKoinu(5000)

	_, err := client.CreateTransfer(context.Background(), connect.NewRequest(msg))
	requireConnectError(t, err, connect.CodeUnauthenticated)
}

func TestCreateTransfer_MissingRecipientEmail(t *testing.T) {
	client := newTestClient(t)
	token := loginAs(t, client, "void@example.com")

	msg := new(demov1.CreateTransferRequest)
	msg.SetAmountKoinu(5000)

	_, err := client.CreateTransfer(context.Background(), withAuth(msg, token))
	requireConnectError(t, err, connect.CodeInvalidArgument)
}

func TestCreateTransfer_InvalidRecipientEmail(t *testing.T) {
	client := newTestClient(t)
	token := loginAs(t, client, "void@example.com")

	msg := new(demov1.CreateTransferRequest)
	msg.SetRecipientEmail("not-valid")
	msg.SetAmountKoinu(5000)

	_, err := client.CreateTransfer(context.Background(), withAuth(msg, token))
	requireConnectError(t, err, connect.CodeInvalidArgument)
}

func TestCreateTransfer_AmountAtDustLimit(t *testing.T) {
	client := newTestClient(t)
	token := loginAs(t, client, "void@example.com")

	msg := new(demov1.CreateTransferRequest)
	msg.SetRecipientEmail("shiba@example.com")
	msg.SetAmountKoinu(1000) // proto constraint is gt 1000, so 1000 is rejected

	_, err := client.CreateTransfer(context.Background(), withAuth(msg, token))
	requireConnectError(t, err, connect.CodeInvalidArgument)
}

func TestCreateTransfer_InsufficientBalance(t *testing.T) {
	client := newTestClient(t)
	token := loginAs(t, client, "void@example.com")

	msg := new(demov1.CreateTransferRequest)
	msg.SetRecipientEmail("shiba@example.com")
	msg.SetAmountKoinu(9_999_999_999)

	_, err := client.CreateTransfer(context.Background(), withAuth(msg, token))
	requireConnectError(t, err, connect.CodeFailedPrecondition)
}
