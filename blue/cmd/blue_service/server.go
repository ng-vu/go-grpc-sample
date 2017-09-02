package main

import (
	"io/ioutil"
	"net"
	"net/http"
	"path/filepath"

	"github.com/facebookgo/inject"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/ng-vu/go-grpc-sample/base/l"
	"github.com/ng-vu/go-grpc-sample/blue/config"
	"github.com/ng-vu/go-grpc-sample/blue/internal/auth"
	"github.com/ng-vu/go-grpc-sample/blue/internal/controller"
	"github.com/ng-vu/go-grpc-sample/blue/internal/grpc"
	"github.com/ng-vu/go-grpc-sample/blue/internal/service"
	"github.com/ng-vu/go-grpc-sample/blue/internal/store"
	pbAgency "github.com/ng-vu/go-grpc-sample/pb/agency"
	pbPartner "github.com/ng-vu/go-grpc-sample/pb/partner"
	pbSAdmin "github.com/ng-vu/go-grpc-sample/pb/sadmin"
)

type gatewayFunc func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error

func startServer(
	name string,
	GRPC config.GRPC,
	HTTP config.HTTP,
	authFunc grpcTransport.AuthFunc,
	gatewayFunc gatewayFunc,
	configFunc func(*grpc.Server, *runtime.ServeMux) (http.Handler, error),
) (*grpc.Server, *http.Server) {

	logger := l.New()

	// Config GRPC server
	grpcSvr := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpcTransport.LogUnaryServerInterceptor(logger),
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_opentracing.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			grpcTransport.AuthUnaryServerInterceptor(authFunc),
		)),
	)
	ln, err := net.Listen("tcp", GRPC.Listen())
	if err != nil {
		ll.Fatal(name+" GRPC Server error", l.Error(err))
	}

	// Config JSON Gateway
	mux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{OrigName: true, EmitDefaults: true}))
	opts := []grpc.DialOption{grpc.WithInsecure()}
	endpoint := "127.0.0.1:" + GRPC.Port
	err = gatewayFunc(ctx, mux, endpoint, opts)
	if err != nil {
		ll.Fatal(name+" JSON Gateway Partner error", l.Error(err))
	}

	// Register service logic
	handler, err := configFunc(grpcSvr, mux)
	if err != nil {
		ll.Fatal(name+" HTTP Server error", l.Error(err))
	}

	// Start GRPC server
	go func() {
		defer ctxCancel()

		ll.Info(name+" GRPC Server started", l.String("listen", GRPC.Listen()))
		err = grpcSvr.Serve(ln)
		if err != nil {
			ll.Error(name+" GRPC Server error", l.Error(err))
		}
	}()

	// Start HTTP server
	httpSvr := &http.Server{Addr: HTTP.Listen(), Handler: handler}
	go func() {
		defer ctxCancel()

		ll.Info(name+" HTTP server started", l.String("listen", HTTP.Listen()))
		err = httpSvr.ListenAndServe()
		if err != http.ErrServerClosed {
			ll.Error(name+" JSON Gateway error", l.Error(err))
		}
	}()

	return grpcSvr, httpSvr
}

func docHandler(docFile string) http.HandlerFunc {
	docPath := filepath.Join(cfg.DocumentPath, docFile)
	doc, err := ioutil.ReadFile(docPath)
	if err != nil {
		ll.Warn("Unable to read file", l.String("path", docPath))
		ll.Warn("Document is not supported")
		doc = []byte("Document not found")
	}
	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(doc)
	}
}

func populate(objects ...interface{}) {
	err := inject.Populate(objects...)
	if err != nil {
		ll.Fatal("Unable to initialize objects", l.Error(err))
	}
}

func startServers() {
	redisStore := connectRedis(cfg)
	tokenGenerator := auth.NewGenerator("", redisStore)
	db := connectPostgres(cfg)

	agencyStaffStore := store.NewAgencyStaffStore(db)
	orderStore := store.NewOrderStore(db)
	orderTransactionStore := store.NewOrderTransactionStore(db)
	serviceProviderStore := store.NewServiceProviderStore(db)
	serviceStore := store.NewServiceStore(db)
	userInternalStore := store.NewUserInternalStore(db)

	orderCtrl := new(ctrl.OrderCtrl)
	serviceCtrl := new(ctrl.ProviderCtrl)
	userCtrl := ctrl.NewUserCtrl(redisStore, tokenGenerator)

	inner := service.NewInnerService()
	populate(
		redisStore,
		agencyStaffStore,
		orderStore,
		orderTransactionStore,
		serviceProviderStore,
		serviceStore,
		userInternalStore,
		orderCtrl,
		serviceCtrl,
		userCtrl,
		&inner,
	)

	// Initialize config and data
	if err := inner.SetupServices(); err != nil {
		ll.Fatal("Unable to setup services", l.Error(err))
	}

	{
		s := service.NewAgencyService(inner)
		authFunc := grpcTransport.Authentication(tokenGenerator, "",
			[]string{
				"/agency.BlueAgency/AccountLogin",
				"/agency.BlueAgency/VersionInfo",
			})

		grpc1, http1 = startServer("BlueAgency ", cfg.AgencyService.GRPC, cfg.AgencyService.HTTP,
			authFunc, pbAgency.RegisterBlueAgencyHandlerFromEndpoint,
			func(grpcSvr *grpc.Server, api *runtime.ServeMux) (http.Handler, error) {
				pbAgency.RegisterBlueAgencyServer(grpcSvr, s)
				m := http.NewServeMux()
				m.Handle("/", http.RedirectHandler("/doc", http.StatusTemporaryRedirect))
				m.Handle("/doc", http.RedirectHandler("/doc/swagger", http.StatusTemporaryRedirect))
				m.Handle("/doc/swagger", docHandler("agency.swagger.json"))
				m.Handle("/api/", api)
				return m, nil
			})
	}
	{
		s := service.NewSAdminService(inner)
		authFunc := grpcTransport.Authentication(tokenGenerator, cfg.SAdminService.MagicToken,
			[]string{
				"/sadmin.BlueSAdmin/VersionInfo",
			})
		grpc2, http2 = startServer("BlueSAdmin ", cfg.SAdminService.GRPC, cfg.SAdminService.HTTP,
			authFunc, pbSAdmin.RegisterBlueSAdminHandlerFromEndpoint,
			func(grpcSvr *grpc.Server, api *runtime.ServeMux) (http.Handler, error) {
				pbSAdmin.RegisterBlueSAdminServer(grpcSvr, s)
				m := http.NewServeMux()
				m.Handle("/", http.RedirectHandler("/doc", http.StatusTemporaryRedirect))
				m.Handle("/doc", http.RedirectHandler("/doc/swagger", http.StatusTemporaryRedirect))
				m.Handle("/doc/swagger", docHandler("sadmin.swagger.json"))
				m.Handle("/api/", api)
				return m, nil
			})
	}
	{
		s := service.NewPartnerService(inner)
		validator := s.(interface {
			Validate(tokenStr string) (auth.ServiceProviderClaim, error)
		})
		authFunc := func(ctx context.Context, fullMethod string) (context.Context, error) {
			if fullMethod == "/partner.BluePartner/VersionInfo" {
				return ctx, nil
			}

			apiKey, err := grpc_auth.AuthFromMD(ctx, "bearer")
			if err != nil {
				ll.Warn("No authorization header", l.String("method", fullMethod), l.Error(err))
				return ctx, err
			}

			spClaim, err := validator.Validate(apiKey)
			if err != nil {
				ll.Warn("Invalid apiKey", l.String("apiKey", apiKey), l.Error(err))
				return ctx, err
			}

			newCtx := auth.NewContextWithProvider(ctx, spClaim)
			return newCtx, nil
		}

		grpc2, http2 = startServer("BluePartner", cfg.PartnerService.GRPC, cfg.PartnerService.HTTP,
			authFunc, pbPartner.RegisterBluePartnerHandlerFromEndpoint,
			func(grpcSvr *grpc.Server, api *runtime.ServeMux) (http.Handler, error) {
				pbPartner.RegisterBluePartnerServer(grpcSvr, s)
				m := http.NewServeMux()
				m.Handle("/", http.RedirectHandler("/doc", http.StatusTemporaryRedirect))
				m.Handle("/doc", http.RedirectHandler("/doc/swagger", http.StatusTemporaryRedirect))
				m.Handle("/doc/swagger", docHandler("partner.swagger.json"))
				m.Handle("/api/", api)
				return m, nil
			})
	}
}
