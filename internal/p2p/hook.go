package p2p

type Hook interface {
    OnPeerJoin(id string)
    OnPeerLeave(id string)
}

type NopHook struct{}
func (NopHook) OnPeerJoin(id string) {}
func (NopHook) OnPeerLeave(id string) {}