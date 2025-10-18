package main

import (
    "context"
    "flag"
    "os"
    "os/signal"
    "syscall"

    "modular-dvt-engine/internal/api"
    "modular-dvt-engine/internal/consensus"
    "modular-dvt-engine/internal/monitoring"
    "modular-dvt-engine/internal/p2p"
    "modular-dvt-engine/pkg/lifecycle"
    "modular-dvt-engine/pkg/logger"
)

func main() {
    var (
        apiAddr  string
        monAddr  string
    )
    flag.StringVar(&apiAddr, "validator-api", "127.0.0.1:4600", "Validator API listen address")
    flag.StringVar(&monAddr, "monitoring", "127.0.0.1:4620", "Monitoring listen address")
    flag.Parse()

    ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
    defer cancel()

    m := lifecycle.New()
    m.Add(api.New(apiAddr))
    m.Add(monitoring.New(monAddr))
    m.Add(p2p.New())
    m.Add(consensus.New())

    if err := m.StartAll(ctx); err != nil { logger.Error(err.Error()); os.Exit(1) }
    <-ctx.Done()
    _ = m.StopAll(context.Background())
}
