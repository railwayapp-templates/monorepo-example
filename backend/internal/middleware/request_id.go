package middleware

import (
	"context"
	"crypto/rand"
	"fmt"
	"net/http"
)

// RequestIDKey is the key that holds the unique request ID in a request context.
var requestIDKey = struct{}{}

// RequestID is a middleware that injects a request ID into the context of each request.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), requestIDKey, createID(14))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetReqID returns a request ID from the given context if one is present.
func GetReqID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	if reqID, ok := ctx.Value(requestIDKey).(string); ok {
		return reqID
	}

	return ""
}

func createID(len int) string {
	b := make([]byte, len)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}

	return fmt.Sprintf("%04x", b)
}
