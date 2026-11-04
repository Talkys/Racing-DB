package routes

import (
	"context"
	"log/slog"
	"net/http"
	"regexp"
	"runtime/debug"
	"time"

	"github.com/google/uuid"
)

const RequestIDHeader = "X-Request-Id"

// An incoming id is reused only if it is short and harmless to log; otherwise a new one is generated.
var requestIDPattern = regexp.MustCompile(`^[A-Za-z0-9._-]{1,64}$`)

type requestIDKey struct{}

// RequestID returns the id the Middleware gave this request, or "" outside it.
func RequestID(ctx context.Context) string {
	id, _ := ctx.Value(requestIDKey{}).(string)
	return id
}

// statusRecorder remembers the status a handler wrote, so the access log can report it.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (rec *statusRecorder) WriteHeader(status int) {
	if rec.status == 0 {
		rec.status = status
	}
	rec.ResponseWriter.WriteHeader(status)
}

func (rec *statusRecorder) Write(b []byte) (int, error) {
	if rec.status == 0 {
		rec.status = http.StatusOK
	}
	return rec.ResponseWriter.Write(b)
}

// Unwrap lets http.ResponseController reach the underlying writer.
func (rec *statusRecorder) Unwrap() http.ResponseWriter {
	return rec.ResponseWriter
}

// Middleware gives each request an id (X-Request-Id), recovers panics as a 500 and writes one access-log line.
func Middleware(log *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		id := r.Header.Get(RequestIDHeader)
		if !requestIDPattern.MatchString(id) {
			id = uuid.NewString()
		}
		w.Header().Set(RequestIDHeader, id)
		r = r.WithContext(context.WithValue(r.Context(), requestIDKey{}, id))

		rec := &statusRecorder{ResponseWriter: w}

		defer func() {
			if v := recover(); v != nil {
				// ErrAbortHandler is net/http's way to abort a response on purpose; let the server handle it.
				if v == http.ErrAbortHandler {
					panic(v)
				}
				log.Error("panic", "id", id, "panic", v, "stack", string(debug.Stack()))
				if rec.status == 0 {
					WriteJSONResponse(rec, http.StatusInternalServerError,
						ErrorResponse{Status: http.StatusInternalServerError, Error: "internal server error"})
				}
			}

			status := rec.status
			if status == 0 {
				status = http.StatusOK
			}
			log.Info("request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", status,
				"latency", time.Since(start),
				"id", id,
			)
		}()

		next.ServeHTTP(rec, r)
	})
}
