Aequa Network — Modular DVT Engine

A minimal, production‑minded Distributed Validator (DVT) engine with strong observability and security gates. It provides a stable, layered foundation to plug in privacy encryption, batch auctions (DFBA) and external solvers later — without breaking metrics/logging dimensions.

Status

- M3 complete (foundation usable):
  - API: `/v1/duty` strict validation (≤1MiB; type/height/round), unified JSON logs and Prom summaries
  - P2P: config‑driven gates (AllowList→Rate→Score) + resource limits (MaxConns); DKG/cluster‑lock precheck (fail‑fast)
  - Consensus: QBFT verifier (strict + anti‑replay) and state skeleton (preprepare/prepare/commit with dedup); StateDB atomic persist + pessimistic recovery
  - Observability: stable fields (trace_id/route/code/result/latency_ms) and metrics families (`api_requests_total`, `service_op_ms`, `consensus_*`, `p2p_*`)
- M4 in progress:
  - Strengthen tests (fuzz/adversarial/chaos) on frozen interfaces and dimensions
  - Extend `consensus_proc_ms` timing to cover verify→state→persist end‑to‑end (included in this update)

Quick Start

```bash
# Build
go build -o bin/dvt-node ./cmd/dvt-node

# Run (local)
./bin/dvt-node --validator-api 127.0.0.1:4600 --monitoring 127.0.0.1:4620

# Health
curl http://127.0.0.1:4600/health  # -> ok

# Metrics
curl http://127.0.0.1:4620/metrics
```

Or via Docker:

```bash
docker build -t aequa-local:latest .
docker compose up -d
```

Observability (Stable)

- Logs (JSON):
  - `api_request` (route, code, latency_ms, result, trace_id, err?)
  - `service_op` (service, op, latency_ms, result, err?)
  - `consensus_recv` (kind, trace_id, latency_ms)
  - `qbft_verify`, `qbft_state`, `p2p_peer`, `consensus_state`
- Metrics (Prometheus):
  - `api_requests_total{route,code}`, `api_latency_ms_sum/_count{route}`
  - `service_op_ms_sum/_count{service,op}`
  - `consensus_events_total{kind}`, `consensus_proc_ms_sum/_count{kind}`
  - `qbft_msg_verified_total{result|type}`, `qbft_state_transitions_total{type}`
  - `p2p_conn_attempts_total{result}`, `p2p_conns_open`, `p2p_conn_open_total/close_total`
  - `state_persist_ms_sum/_count`, `state_recovery_total{result}`

CI / Security Gates

- Go 1.24; `golangci-lint`; unit tests with coverage gate (≥70% on internal/api)
- `govulncheck`; Snyk CLI (token required); UTF‑8 no‑BOM check (blocking)
- QBFT tests run as part of required checks

License

Business Source License 1.1 (BSL 1.1). See LICENSE for details.
