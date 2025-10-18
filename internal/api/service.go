package api

import (
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/http/httputil"
    "net/url"
    "time"

    "github.com/zimingliu11111111/Aequa-network/pkg/lifecycle"
    "github.com/zimingliu11111111/Aequa-network/pkg/logger"
    "github.com/zimingliu11111111/Aequa-network/pkg/metrics"
)

type Service struct{ addr string; srv *http.Server; onPublish func(ctx context.Context, payload []byte) error; upstream string }

func New(addr string, onPublish func(ctx context.Context, payload []byte) error, upstream string) *Service {
    return &Service{addr: addr, onPublish: onPublish, upstream: upstream}
}

func (s *Service) Name() string { return "api" }

func (s *Service) Start(ctx context.Context) error {
    mux := http.NewServeMux()
    mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(200); _, _ = w.Write([]byte("ok")) })
    mux.HandleFunc("/v1/duty", s.handleDuty)
    mux.HandleFunc("/", s.proxy)
    s.srv = &http.Server{ Addr: s.addr, Handler: mux }
    go func() {
        logger.Info(fmt.Sprintf("api listening on %s\n", s.addr))
        if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            logger.Error("api server error: "+err.Error())
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

// handleDuty accepts a JSON body and publishes it after basic validation.
func (s *Service) handleDuty(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost { http.Error(w, "method not allowed", http.StatusMethodNotAllowed); return }
    if r.Body == nil { http.Error(w, "empty body", http.StatusBadRequest); return }
    b, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 1<<20))
    if err != nil { http.Error(w, "read error", http.StatusBadRequest); return }
    if err := validateDutyJSON(b); err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }

    start := time.Now()
    if s.onPublish != nil { _ = s.onPublish(r.Context(), b) }
    dur := time.Since(start)
    metrics.Inc("api_requests_total", map[string]string{"route":"/v1/duty","code":"202"})
    logger.Info(fmt.Sprintf("duty_accept route=/v1/duty bytes=%d duration_ms=%d", len(b), dur.Milliseconds()))
    w.WriteHeader(http.StatusAccepted)
}

type dutyEnvelope struct {
    Type   string `json:"type"`
    Height uint64 `json:"height"`
    Round  uint64 `json:"round"`
    Payload any   `json:"payload"`
}

func validateDutyJSON(b []byte) error {
    var d dutyEnvelope
    if len(b) == 0 { return fmt.Errorf("empty") }
    if len(b) > 1<<20 { return fmt.Errorf("too large") }
    if err := json.Unmarshal(b, &d); err != nil { return fmt.Errorf("invalid json") }
    switch d.Type {
    case "attester", "proposer", "sync":
    default:
        return fmt.Errorf("invalid type")
    }
    if d.Height > 1<<62 { return fmt.Errorf("height out of range") }
    if d.Round > 1<<40 { return fmt.Errorf("round out of range") }
    return nil
}

func (s *Service) proxy(w http.ResponseWriter, r *http.Request) {
    if s.upstream == "" { http.NotFound(w, r); return }
    u, err := url.Parse(s.upstream)
    if err != nil { http.Error(w, "bad upstream", http.StatusBadGateway); return }
    rp := httputil.NewSingleHostReverseProxy(u)
    start := time.Now()
    rp.ModifyResponse = func(resp *http.Response) error {
        metrics.Inc("api_requests_total", map[string]string{"route":"proxy","code":fmt.Sprintf("%d", resp.StatusCode)})
        logger.InfoJ("proxy", map[string]any{"code": resp.StatusCode, "duration_ms": time.Since(start).Milliseconds()})
        return nil
    }
    rp.ErrorHandler = func(w http.ResponseWriter, r *http.Request, e error) {
        metrics.Inc("api_requests_total", map[string]string{"route":"proxy","code":"502"})
        logger.ErrorJ("proxy_error", map[string]any{"err": e.Error()})
        http.Error(w, "upstream error", http.StatusBadGateway)
    }
    rp.ServeHTTP(w, r)
}



