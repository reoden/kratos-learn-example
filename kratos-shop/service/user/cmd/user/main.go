package main

import (
    "flag"
    "os"

    "github.com/go-kratos/kratos/v2"
    "github.com/go-kratos/kratos/v2/config"
    "github.com/go-kratos/kratos/v2/config/file"
    "github.com/go-kratos/kratos/v2/log"
    "github.com/go-kratos/kratos/v2/middleware/tracing"
    "github.com/go-kratos/kratos/v2/registry"
    "github.com/go-kratos/kratos/v2/transport/grpc"
    "user/internal/conf"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
    // Name is the name of the compiled software.
    Name = "shop.users.service"
    // Version is the version of the compiled software.
    Version = "v1"
    // flagconf is the config flag.
    flagconf string

    id, _ = os.Hostname()
)

func init() {
    flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func newApp(logger log.Logger, gs *grpc.Server, rr registry.Registrar) *kratos.App {
    return kratos.New(
        kratos.ID(id+"shop.user.service"),
        kratos.Name(Name),
        kratos.Version(Version),
        kratos.Metadata(map[string]string{}),
        kratos.Logger(logger),
        kratos.Server(
            gs,
        ),
        kratos.Registrar(rr), // consul 的引入
    )
}

func main() {
    flag.Parse()
    logger := log.With(log.NewStdLogger(os.Stdout),
        "ts", log.DefaultTimestamp,
        "caller", log.DefaultCaller,
        "service.id", id,
        "service.name", Name,
        "service.version", Version,
        "trace_id", tracing.TraceID(),
        "span_id", tracing.SpanID(),
    )
    c := config.New(
        config.WithSource(
            file.NewSource(flagconf),
        ),
    )
    defer c.Close()

    if err := c.Load(); err != nil {
        panic(err)
    }

    var bc conf.Bootstrap
    if err := c.Scan(&bc); err != nil {
        panic(err)
    }
    // consul 的引入
    var rc conf.Registry 
    if err := c.Scan(&rc); err != nil {
        panic(err)
    }
    app, cleanup, err := wireApp(bc.Server, bc.Data, &rc, logger)
    if err != nil {
        panic(err)
    }
    defer cleanup()

    // start and wait for stop signal
    if err := app.Run(); err != nil {
        panic(err)
    }
}
