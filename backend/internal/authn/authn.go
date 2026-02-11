package authn

import (
	"context"
	"fmt"
	"log/slog"

	"connectrpc.com/connect"
	"example.com/internal/service"
)

type UnaryInterceptorOptions struct {
	UnauthenticatedProcedures map[string]struct{}
	GetUserFunc               func(ctx context.Context, authToken service.AuthToken) (*service.User, error)
}

func UnaryInterceptor(opts UnaryInterceptorOptions) (connect.UnaryInterceptorFunc, error) {
	if opts.UnauthenticatedProcedures == nil {
		return nil, fmt.Errorf("opts.UnauthenticatedProcedures is nil")
	}
	if opts.GetUserFunc == nil {
		return nil, fmt.Errorf("opts.GetUserFunc is nil")
	}

	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			procedure := req.Spec().Procedure
			_, procDoesNotRequireAuth := opts.UnauthenticatedProcedures[procedure]
			procRequiresAuth := !procDoesNotRequireAuth

			authToken := req.Header().Get("Demo-Auth-Token")
			slog.Info("authn.UnaryInterceptor",
				"procedure", procedure,
				"procDoesNotRequireAuth", procDoesNotRequireAuth,
				"authToken", authToken,
			)

			if authToken == "" {
				if procRequiresAuth {
					return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("missing Demo-Auth-Token header"))
				}
				return next(ctx, req)
			}

			user, err := opts.GetUserFunc(ctx, service.AuthToken(authToken))
			if err != nil || user == nil || user.Email == "" {
				return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("authentication required"))
			}

			ctx = NewContext(ctx, user)

			return next(ctx, req)
		}
	}, nil
}
