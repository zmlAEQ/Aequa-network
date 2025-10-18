package monitoring

import (
    "context"
    "fmt"
    "net/http"
    "time"

    "github.com/zimingliu11111111/Aequa-network/pkg/lifecycle"
    "github.com/zimingliu11111111/Aequa-network/pkg/logger"
)

type Service struct{ addr string; srv *http.Server }

func New(addr string) *Service { return &Service{addr: addr} }
func (s *Service) Name() string { return "monitoring" }

func (s *Service) Start(ctx context.Context) error {
    mux := http.NewServeMux()
    mux.HandleFunc("/metrics", func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte("# HELP dvt_up 1\n# TYPE dvt_up gauge\ndvt_up 1\n")) })
    s.srv = &http.Server{ Addr: s.addr, Handler: mux }
    go func() {
        logger.Info(fmt.Sprintf("monitoring on %s\n", s.addr))
        if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            logger.Error("monitoring server error: "+err.Error())
        }
    }()
    return nil
}

func (s *Service) Stop(ctx context.Context) error {
    if s.srv == nil { return nil }
    ctx, cancel := context.WithTimeout(ctx, 3*time.Second); defer cancel()
    return s.srv.Shutdown(ctx)
}

var _ lifecycle.Service = (*Service)(nil)

