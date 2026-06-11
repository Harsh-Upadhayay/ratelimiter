package ratelimiter

type record struct {
	state   algorithmState
	version int
}

type StateStore interface {
	Get(key string) (algorithmState, int, bool, error)
	CompareAndSwap(key string, version int, state algorithmState) (bool, error)
}
