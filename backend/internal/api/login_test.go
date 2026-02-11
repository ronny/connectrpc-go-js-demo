package api_test

import (
	"context"
	"strings"
	"testing"

	"connectrpc.com/connect"
	demov1 "example.com/gen/myorg/demo/v1"
)

func TestLogin_Success(t *testing.T) {
	client := newTestClient(t)

	req := new(demov1.LoginRequest)
	req.SetEmail("void@example.com")
	req.SetPassword("dogecoin")

	resp, err := client.Login(context.Background(), connect.NewRequest(req))
	if err != nil {
		t.Fatal(err)
	}

	if len(resp.Msg.GetAuthToken()) < 16 {
		t.Fatalf("expected auth token >= 16 chars, got %d", len(resp.Msg.GetAuthToken()))
	}
}

func TestLogin_MissingEmail(t *testing.T) {
	client := newTestClient(t)

	req := new(demov1.LoginRequest)
	req.SetPassword("dogecoin")

	_, err := client.Login(context.Background(), connect.NewRequest(req))
	requireConnectError(t, err, connect.CodeInvalidArgument)
}

func TestLogin_InvalidEmailFormat(t *testing.T) {
	client := newTestClient(t)

	req := new(demov1.LoginRequest)
	req.SetEmail("not-an-email")
	req.SetPassword("dogecoin")

	_, err := client.Login(context.Background(), connect.NewRequest(req))
	requireConnectError(t, err, connect.CodeInvalidArgument)
}

func TestLogin_EmailTooLong(t *testing.T) {
	client := newTestClient(t)

	req := new(demov1.LoginRequest)
	req.SetEmail(strings.Repeat("a", 191) + "@example.com") // 203 chars
	req.SetPassword("dogecoin")

	_, err := client.Login(context.Background(), connect.NewRequest(req))
	requireConnectError(t, err, connect.CodeInvalidArgument)
}

func TestLogin_MissingPassword(t *testing.T) {
	client := newTestClient(t)

	req := new(demov1.LoginRequest)
	req.SetEmail("void@example.com")

	_, err := client.Login(context.Background(), connect.NewRequest(req))
	requireConnectError(t, err, connect.CodeInvalidArgument)
}

func TestLogin_WrongPassword(t *testing.T) {
	client := newTestClient(t)

	req := new(demov1.LoginRequest)
	req.SetEmail("void@example.com")
	req.SetPassword("bitcoin")

	_, err := client.Login(context.Background(), connect.NewRequest(req))
	requireConnectError(t, err, connect.CodeInvalidArgument)
}
