package monitoring

import (
    "context"
    "fmt"
    "net/http"
    "time"

    "github.com/zmlAEQ/Aequa-network/pkg/lifecycle"
    "github.com/zmlAEQ/Aequa-network/pkg/logger"
    "github.com/zmlAEQ/Aequa-network/pkg/metrics"
)

type Service struct{ addr string; srv *http.Server }

func New(addr string) *Service { return &Service{addr: addr} }
func (s *Service) Name() string { return "monitoring" }

func (s *Service) Start(ctx context.Context) error {
    mux := http.NewServeMux()
    mux.HandleFunc("/metrics", s.handleMetrics)
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

// handleMetrics returns current Prom exposition and records unified logs + summary.
func (s *Service) handleMetrics(w http.ResponseWriter, r *http.Request) {
    start := time.Now()
    body := metrics.DumpProm()
    _, _ = w.Write([]byte(body))
    dur := time.Since(start)
    metrics.Inc("api_requests_total", map[string]string{"route":"/metrics","code":"200"})
    metrics.ObserveSummary("api_latency_ms", map[string]string{"route":"/metrics"}, float64(dur.Milliseconds()))
    logger.InfoJ("api_request", map[string]any{
        "route": "/metrics",
        "code": 200,
        "latency_ms": dur.Milliseconds(),
        "result": "ok",
        "trace_id": traceID(r),
    })
}

// traceID returns request trace id from header or generates a simple one.
func traceID(r *http.Request) string {
    if t := r.Header.Get("X-Trace-ID"); t != "" { return t }
    return fmt.Sprintf("%d", time.Now().UnixNano())
}

