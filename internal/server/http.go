package server

import (
	"net/http"

	"github.com/cyc1ones/go-kit/errors"
	"github.com/cyc1ones/go-kit/kratos/middleware/realip"
	v1 "github.com/go-kratos/kratos-layout/api/helloworld/v1"
	versionv1 "github.com/go-kratos/kratos-layout/api/version/v1"
	"github.com/go-kratos/kratos-layout/internal/conf"
	"github.com/go-kratos/kratos-layout/internal/service"

	"github.com/go-kratos/kratos/contrib/middleware/validate/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
)

func errorEncoder(w http.ResponseWriter, r *http.Request, err error) {
	khttp.DefaultErrorEncoder(w, r, errors.Sanitize(err))
}

func newRealIPMiddleware(conf *conf.Server_HTTP_RealIP, logger log.Logger) middleware.Middleware {
	return realip.New(
		realip.WithLogger(logger),
		realip.WithIPHeaders(conf.IpHeaders...),
		realip.WithProxies(conf.TrustedProxies),
	).Server()
}

// NewHTTPServer new an HTTP server.
func NewHTTPServer(
	c *conf.Server,
	logger log.Logger,
	greeter *service.GreeterService,
	version *service.VersionService,
) *khttp.Server {
	// middlewares ...
	var middlewares = []middleware.Middleware{
		tracing.Server(),
		logging.Server(logger),
		recovery.Recovery(),
		metadata.Server(),
	}
	if c.Http.RealIp != nil {
		middlewares = append(middlewares, newRealIPMiddleware(c.Http.RealIp, logger))
	}
	middlewares = append(middlewares, validate.ProtoValidate())

	// options...
	var opts = []khttp.ServerOption{
		khttp.Middleware(middlewares...),
		khttp.ErrorEncoder(errorEncoder),
		khttp.Logger(logger),
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

	// register services...
	srv := khttp.NewServer(opts...)
	v1.RegisterGreeterHTTPServer(srv, greeter)
	versionv1.RegisterVersionHTTPServer(srv, version)
	return srv
}
