package api_test

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	demov1 "example.com/gen/myorg/demo/v1"
)

func TestGetBalance_VoidUser(t *testing.T) {
	client := newTestClient(t)
	token := loginAs(t, client, "void@example.com")

	req := withAuth(new(demov1.GetBalanceRequest), token)

	resp, err := client.GetBalance(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Msg.GetKoinu() != 1_234_567_890 {
		t.Fatalf("expected 1234567890, got %d", resp.Msg.GetKoinu())
	}
}

func TestGetBalance_ShibaUser(t *testing.T) {
	client := newTestClient(t)
	token := loginAs(t, client, "shiba@example.com")

	req := withAuth(new(demov1.GetBalanceRequest), token)

	resp, err := client.GetBalance(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Msg.GetKoinu() != 2_111_222_333 {
		t.Fatalf("expected 2111222333, got %d", resp.Msg.GetKoinu())
	}
}

func TestGetBalance_MissingAuthToken(t *testing.T) {
	client := newTestClient(t)

	_, err := client.GetBalance(context.Background(), connect.NewRequest(new(demov1.GetBalanceRequest)))
	requireConnectError(t, err, connect.CodeUnauthenticated)
}

func TestGetBalance_InvalidAuthToken(t *testing.T) {
	client := newTestClient(t)

	req := withAuth(new(demov1.GetBalanceRequest), "bogus-token-value")
	_, err := client.GetBalance(context.Background(), req)
	requireConnectError(t, err, connect.CodeUnauthenticated)
}
