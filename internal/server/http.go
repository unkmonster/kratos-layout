package server

import (
	"net/http"

	v1 "github.com/go-kratos/kratos-layout/api/helloworld/v1"
	"github.com/go-kratos/kratos-layout/internal/conf"
	"github.com/go-kratos/kratos-layout/internal/service"
	"github.com/gorilla/handlers"
	"github.com/unkmonster/go-kit/middleware/http/realip"

	"github.com/go-kratos/kratos/contrib/middleware/validate/v2"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
)

func newCorsHandler(c *conf.Server_HTTP_CORS) func(http.Handler) http.Handler {
	opts := []handlers.CORSOption{
		handlers.AllowedHeaders(c.AllowedHeaders),
		handlers.AllowedMethods(c.AllowedMethods),
		handlers.AllowedOrigins(c.AllowedOrigins),
	}

	if c.AllowCredentials {
		opts = append(opts, handlers.AllowCredentials())
	}

	return handlers.CORS(
		opts...,
	)
}

// hideInternalErrorEncoder 如果 err 不为 kratos error 响应 500 INTERNAL_SERVER_ERROR 给客户端
func hideInternalErrorEncoder(w http.ResponseWriter, r *http.Request, err error) {
	var ke *errors.Error
	if !errors.As(err, &ke) {
		ke = errors.InternalServer(
			"INTERNAL_SERVER_ERROR",
			"internal server error",
		).WithCause(err)
	}
	khttp.DefaultErrorEncoder(w, r, ke)
}

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.GreeterService, logger log.Logger) *khttp.Server {
	// middlewares ...
	var middlewares = []middleware.Middleware{
		recovery.Recovery(),
		tracing.Server(),
		metadata.Server(),
		logging.Server(logger),
	}
	if c.Http.RealIp != nil {
		middlewares = append(middlewares, realip.Server(
			logger,
			realip.WithIpHeaders(c.Http.RealIp.IpHeaders),
			realip.WithTrustedHeader(c.Http.RealIp.TrustedHeader),
			realip.WithTrustedProxies(c.Http.RealIp.TrustedProxies),
		))
	}
	middlewares = append(middlewares, validate.ProtoValidate())

	// filters ...
	var filters []khttp.FilterFunc
	if c.Http.Cors != nil {
		filters = append(filters, newCorsHandler(c.Http.Cors))
	}

	var opts = []khttp.ServerOption{
		khttp.Middleware(middlewares...),
		khttp.Filter(filters...),
		khttp.ErrorEncoder(hideInternalErrorEncoder),
	}
	if c.Http.Network != "" {
		opts = append(opts, khttp.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, khttp.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, khttp.Timeout(c.Http.Timeout.AsDuration()))
	}

	srv := khttp.NewServer(opts...)
	v1.RegisterGreeterHTTPServer(srv, greeter)
	return srv
}
