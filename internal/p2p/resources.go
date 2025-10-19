package p2p

import (
    "sync/atomic"
    "github.com/zmlAEQ/Aequa-network/pkg/metrics"
)

// ResourceLimits define basic resource caps for P2P.
type ResourceLimits struct{ MaxConns int64 }

func DefaultResourceLimits() ResourceLimits { return ResourceLimits{MaxConns: 128} }

// ResourceManager tracks open connections against limits.
type ResourceManager struct{
    limits ResourceLimits
    open   int64
}

func NewResourceManager(l ResourceLimits) *ResourceManager { return &ResourceManager{limits: l} }

// TryOpen increments the open counter if under limit.
func (r *ResourceManager) TryOpen() bool {
    for {
        o := atomic.LoadInt64(&r.open)
        if o >= r.limits.MaxConns { return false }
        if atomic.CompareAndSwapInt64(&r.open, o, o+1) { metrics.Inc("p2p_conn_open_total", nil); return true }
    }
}

// Close decrements the open counter (no-op if already zero).
func (r *ResourceManager) Close() {
    for {
        o := atomic.LoadInt64(&r.open)
        if o <= 0 { return }
        if atomic.CompareAndSwapInt64(&r.open, o, o-1) { metrics.Inc("p2p_conn_close_total", nil); return }
    }
}