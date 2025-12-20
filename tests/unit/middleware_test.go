package unit

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/danmaciel/api/internal/middleware"
	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	// Capture log output
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	loggerMiddleware := middleware.Logger(handler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	loggerMiddleware.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	logOutput := buf.String()
	assert.Contains(t, logOutput, "[GET]")
	assert.Contains(t, logOutput, "/test")
	assert.Contains(t, logOutput, "Status: 200")
	assert.Contains(t, logOutput, "Duration:")
}

func TestLogger_WithCustomStatusCode(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
	})

	loggerMiddleware := middleware.Logger(handler)

	req := httptest.NewRequest(http.MethodPost, "/api/test", nil)
	rec := httptest.NewRecorder()

	loggerMiddleware.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)

	logOutput := buf.String()
	assert.Contains(t, logOutput, "[POST]")
	assert.Contains(t, logOutput, "/api/test")
	assert.Contains(t, logOutput, "Status: 404")
}

func TestRecovery(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	recoveryMiddleware := middleware.Recovery(handler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	recoveryMiddleware.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Internal Server Error")

	logOutput := buf.String()
	assert.Contains(t, logOutput, "PANIC: test panic")
}

func TestRecovery_NoPanic(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	recoveryMiddleware := middleware.Recovery(handler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	recoveryMiddleware.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "OK", rec.Body.String())
}

func TestContentType(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	contentTypeMiddleware := middleware.ContentType("application/json")(handler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	contentTypeMiddleware.ServeHTTP(rec, req)

	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.Equal(t, "OK", rec.Body.String())
}

func TestContentType_CustomType(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<html></html>"))
	})

	contentTypeMiddleware := middleware.ContentType("text/html")(handler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	contentTypeMiddleware.ServeHTTP(rec, req)

	assert.Equal(t, "text/html", rec.Header().Get("Content-Type"))
}

func TestResponseWriter_WriteHeader(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Created"))
	})

	loggerMiddleware := middleware.Logger(handler)

	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	rec := httptest.NewRecorder()

	loggerMiddleware.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestMiddlewareChain(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Chain multiple middlewares
	wrapped := middleware.ContentType("application/json")(
		middleware.Logger(
			middleware.Recovery(handler),
		),
	)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	wrapped.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.Contains(t, buf.String(), "[GET]")
}

func TestMiddlewareChain_WithPanic(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("middleware chain panic")
	})

	wrapped := middleware.ContentType("application/json")(
		middleware.Logger(
			middleware.Recovery(handler),
		),
	)

	req := httptest.NewRequest(http.MethodPost, "/api/test", nil)
	rec := httptest.NewRecorder()

	wrapped.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	// http.Error sets Content-Type to "text/plain; charset=utf-8"
	assert.Equal(t, "text/plain; charset=utf-8", rec.Header().Get("Content-Type"))

	logOutput := buf.String()
	assert.Contains(t, logOutput, "PANIC: middleware chain panic")
}

func TestRecovery_DifferentPanicTypes(t *testing.T) {
	tests := []struct {
		name      string
		panicVal  interface{}
		expectLog string
	}{
		{
			name:      "string panic",
			panicVal:  "string error",
			expectLog: "PANIC: string error",
		},
		{
			name:      "int panic",
			panicVal:  123,
			expectLog: "PANIC: 123",
		},
		{
			name:      "struct panic",
			panicVal:  struct{ msg string }{"struct error"},
			expectLog: "PANIC:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer log.SetOutput(os.Stderr)

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic(tt.panicVal)
			})

			recoveryMiddleware := middleware.Recovery(handler)

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			rec := httptest.NewRecorder()

			recoveryMiddleware.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusInternalServerError, rec.Code)
			logOutput := buf.String()
			assert.True(t, strings.Contains(logOutput, tt.expectLog))
		})
	}
}
