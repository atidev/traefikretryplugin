package traefikretryplugin

import (
	"net/http"
	"strconv"
)

func NewRetryResponseWriter(rw http.ResponseWriter, policy *RetryPolicy, attempt int) *RetryResponseWriter {
	return &RetryResponseWriter{
		rw:      rw,
		policy:  policy,
		attempt: attempt,
	}
}

type RetryResponseWriter struct {
	rw       http.ResponseWriter
	policy   *RetryPolicy
	Retrying bool
	writing  bool
	attempt  int
}

func (w *RetryResponseWriter) shouldRetry(status int) bool {
	return w.policy != nil && w.policy.Applicable(status) && w.policy.CanRetry(w.attempt)
}

func (w *RetryResponseWriter) Header() http.Header {
	if !w.writing {
		return make(http.Header)
	}

	return w.rw.Header()
}

func (w *RetryResponseWriter) WriteHeader(status int) {
	if w.shouldRetry(status) {
		w.Retrying = true
		return
	}

	w.writing = true

	h := w.Header()

	if w.attempt > 0 {
		h.Add("Retry-Attempt", strconv.Itoa(w.attempt))
	}

	w.rw.WriteHeader(status)
}

func (w *RetryResponseWriter) Write(body []byte) (int, error) {
	if w.Retrying {
		return len(body), nil
	}

	return w.rw.Write(body)
}
