package main

import (
    "context"
    "flag"
    "os"
    "os/signal"
    "syscall"

    "github.com/zimingliu11111111/Aequa-network/internal/api"
    "github.com/zimingliu11111111/Aequa-network/internal/consensus"
    "github.com/zimingliu11111111/Aequa-network/internal/monitoring"
    "github.com/zimingliu11111111/Aequa-network/internal/p2p"
    "github.com/zimingliu11111111/Aequa-network/pkg/bus"
    "github.com/zimingliu11111111/Aequa-network/pkg/lifecycle"
    "github.com/zimingliu11111111/Aequa-network/pkg/logger"
)

func main() {
    var (
        apiAddr string
        monAddr string
    )
    flag.StringVar(&apiAddr, "validator-api", "127.0.0.1:4600", "Validator API listen address")
    flag.StringVar(&monAddr, "monitoring", "127.0.0.1:4620", "Monitoring listen address")
    flag.Parse()

    ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
    defer cancel()

    b := bus.New(256)
    publish := func(ctx context.Context, payload []byte) error {
        // TODO: decode, map to consensus Event and publish on bus
        b.Publish(ctx, bus.Event{Kind: bus.KindDuty, Body: payload})
        return nil
    }

    m := lifecycle.New()
    m.Add(api.New(apiAddr, publish))
    m.Add(monitoring.New(monAddr))
    m.Add(p2p.New())
    m.Add(consensus.New())

    if err := m.StartAll(ctx); err != nil { logger.Error(err.Error()); os.Exit(1) }
    <-ctx.Done()
    _ = m.StopAll(context.Background())
}
