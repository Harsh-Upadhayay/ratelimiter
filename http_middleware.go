package ratelimiter

import (
	"net/http"
	"strconv"
	"time"
)

type FailurePolicy int

const (
	FailOpen FailurePolicy = iota
	FailClosed
)

type KeyFunc func(*http.Request) string

type MiddlewareOption func(*Middleware)

func WithFailurePolicy(policy FailurePolicy) MiddlewareOption {
	return func(mw *Middleware) {
		mw.failurePolicy = policy
	}
}

type Middleware struct {
	limiter       Limiter
	keyFunc       KeyFunc
	failurePolicy FailurePolicy
}

func durationToSecondsString(d time.Duration) string {
	if d <= 0 {
		return "0"
	}
	return strconv.Itoa(int((d + time.Second - 1) / (time.Second)))
}

func NewMiddleware(limiter Limiter, keyFunc KeyFunc, opts ...MiddlewareOption) (*Middleware, error) {
	if limiter == nil {
		return nil, ErrNilLimiter
	}

	if keyFunc == nil {
		return nil, ErrNilKeyFunc
	}

	mw := &Middleware{
		limiter:       limiter,
		keyFunc:       keyFunc,
		failurePolicy: FailOpen,
	}

	for _, o := range opts {
		o(mw)
	}

	if !(mw.failurePolicy == FailOpen || mw.failurePolicy == FailClosed) {
		return nil, ErrInvalidFailurePolicy
	}
	return mw, nil
}

func (m *Middleware) Wrap(next http.Handler) http.Handler {
	if next == nil {
		panic("ratelimiter: nil next handler")
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := m.keyFunc(r)

		if key == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		result, err := m.limiter.Allow(r.Context(), key)

		if err != nil {
			if m.failurePolicy == FailOpen {
				next.ServeHTTP(w, r)
				return
			}
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(result.Remaining))

		if result.Allowed {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Retry-After", durationToSecondsString(result.RetryAfter))
		w.WriteHeader(http.StatusTooManyRequests)
	})
}
