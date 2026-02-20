package main

import (
	"context"
	"flag"
	"os"

	"github.com/go-kratos/kratos-layout/internal/conf"
	"github.com/go-kratos/kratos-layout/internal/pkg/otel"
	"github.com/go-kratos/kratos-layout/internal/server"
	"github.com/go-kratos/kratos-layout/internal/version"

	"github.com/cyc1ones/go-kit/background"
	"github.com/cyc1ones/go-kit/flag/value"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/env"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	_ "go.uber.org/automaxprocs"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// name is the name of the compiled software.
	// TODO: set this variable
	name string
	// flagconf is the config flag.
	flagconf value.SliceFlag

	id, _ = os.Hostname()
)

func init() {
	flag.Var(&flagconf, "conf", "config path, eg: -conf config.yaml")
}

func newApp(
	logger log.Logger,
	rr registry.Registrar,
	gs *grpc.Server,
	hs *http.Server,
	bg *background.Background,
	scheduler *server.Scheduler,
) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(name),
		kratos.Version(version.Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Registrar(rr),
		kratos.Server(
			gs,
			hs,
			bg,
			scheduler,
		),
	)
}

func main() {
	flag.Parse()

	cs := []config.Source{
		env.NewSource(""),
	}
	for _, path := range flagconf {
		cs = append(cs, file.NewSource(path))
	}

	c := config.New(config.WithSource(cs...))
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	// init global logger
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", name,
		"service.version", version.Version,
		"service.commit", version.Commit,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	logger = log.NewFilter(logger, log.FilterLevel(log.ParseLevel(bc.Observability.Log.Level)))

	// init otel
	res := otel.NewResource(logger, &bc, name)
	close := otel.InitTraceProvider(
		logger,
		otel.NewSampler(&bc),
		otel.NewTraceExporter(logger, bc.Observability),
		res,
	)
	defer close(context.Background())

	app, cleanup, err := wireApp(
		logger,
		&bc.Env,
		bc.Server,
		bc.Data,
		bc.Observability,
		bc.Registry,
		bc.Auth,
	)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
