package api_test

import (
	"context"
	"errors"
	"testing"

	"connectrpc.com/connect"
	demov1 "example.com/gen/myorg/demo/v1"
	"example.com/gen/myorg/demo/v1/demov1connect"
	"example.com/internal/api"
	"example.com/internal/service"
	"go.akshayshah.org/memhttp"
)

func newTestClient(t *testing.T) demov1connect.DemoAPIClient {
	t.Helper()

	svc := service.New()

	router, err := api.NewRouter(svc)
	if err != nil {
		t.Fatal(err)
	}

	srv, err := memhttp.New(router)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := srv.Close(); err != nil {
			t.Error(err)
		}
	})

	return demov1connect.NewDemoAPIClient(srv.Client(), srv.URL())
}

func loginAs(t *testing.T, client demov1connect.DemoAPIClient, email string) string {
	t.Helper()

	req := new(demov1.LoginRequest)
	req.SetEmail(email)
	req.SetPassword("dogecoin")

	resp, err := client.Login(context.Background(), connect.NewRequest(req))
	if err != nil {
		t.Fatalf("loginAs(%q): %v", email, err)
	}

	return resp.Msg.GetAuthToken()
}

func withAuth[T any](msg *T, token string) *connect.Request[T] {
	req := connect.NewRequest(msg)
	req.Header().Set("Demo-Auth-Token", token)

	return req
}

func requireConnectError(t *testing.T, err error, wantCode connect.Code) {
	t.Helper()

	if err == nil {
		t.Fatalf("expected connect.Code %v, got nil error", wantCode)
	}

	var connectErr *connect.Error
	if !errors.As(err, &connectErr) {
		t.Fatalf("expected *connect.Error, got %T: %v", err, err)
	}

	if connectErr.Code() != wantCode {
		t.Fatalf("expected connect.Code %v, got %v: %v", wantCode, connectErr.Code(), connectErr.Message())
	}
}
