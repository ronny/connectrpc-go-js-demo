package api

import (
	"context"
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	"connectrpc.com/grpchealth"
	"connectrpc.com/grpcreflect"
	"connectrpc.com/validate"
	"example.com/gen/myorg/demo/v1/demov1connect"
	"example.com/internal/authn"
	"example.com/internal/httpserv"
	"example.com/internal/service"
)

type Server struct {
	svc        *service.Service
	httpServer *http.Server
}

func NewServer(listenAddr string, svc *service.Service) (*Server, error) {
	protocols := new(http.Protocols)
	protocols.SetHTTP1(true)
	protocols.SetUnencryptedHTTP2(true)

	router := http.NewServeMux()

	compress1KB := connect.WithCompressMinBytes(1024)

	healthCheckProcedure := fmt.Sprintf("/%s/Check", grpchealth.HealthV1ServiceName)
	healthWatchProcedure := fmt.Sprintf("/%s/Watch", grpchealth.HealthV1ServiceName)

	authnInterceptor, err := authn.UnaryInterceptor(authn.UnaryInterceptorOptions{
		UnauthenticatedProcedures: map[string]struct{}{
			healthCheckProcedure:                {},
			healthWatchProcedure:                {},
			demov1connect.DemoAPILoginProcedure: {},
		},
		GetUserFunc: svc.GetUser,
	})
	if err != nil {
		return nil, fmt.Errorf("authn.UnaryInterceptor: %w", err)
	}

	// Order matters
	interceptors := []connect.Interceptor{
		authnInterceptor,
		validate.NewInterceptor(
			validate.WithValidateResponses(),
		),
	}

	handler := NewHandler(svc)

	router.Handle(
		demov1connect.NewDemoAPIHandler(
			handler,
			connect.WithInterceptors(interceptors...),
			compress1KB,
		),
	)

	router.Handle(grpchealth.NewHandler(
		grpchealth.NewStaticChecker(demov1connect.DemoAPIName),
	))

	reflector := grpcreflect.NewStaticReflector(demov1connect.DemoAPIName)
	router.Handle(grpcreflect.NewHandlerV1(reflector, compress1KB))
	router.Handle(grpcreflect.NewHandlerV1Alpha(reflector, compress1KB))

	httpServer := httpserv.WithDefaults(&http.Server{
		Addr:      listenAddr,
		Handler:   router,
		Protocols: protocols,
	})

	return &Server{
		svc:        svc,
		httpServer: httpServer,
	}, nil
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.httpServer.SetKeepAlivesEnabled(false)
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) Address() string {
	return s.httpServer.Addr
}
