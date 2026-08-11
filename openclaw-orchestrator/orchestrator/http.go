package orchestrator

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"
)

type HTTPOptions struct {
	MutationHeader string
	MutationValue  string
}

func NewHTTPHandler(service *Service, options HTTPOptions) http.Handler {
	if options.MutationHeader == "" {
		options.MutationHeader = "X-OpenClaw-Operator-Request"
	}
	if options.MutationValue == "" {
		options.MutationValue = "1"
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if !allowMethod(w, r, http.MethodGet) {
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "service": "openclaw-orchestrator", "time": time.Now().Format(time.RFC3339)})
	})
	mux.HandleFunc("/v1/status", func(w http.ResponseWriter, r *http.Request) {
		if !allowMethod(w, r, http.MethodGet) {
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
		defer cancel()
		status, err := service.Status(ctx)
		writeResult(w, status, err)
	})
	mux.HandleFunc("/v1/models", func(w http.ResponseWriter, r *http.Request) {
		if !allowMethod(w, r, http.MethodGet) {
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 45*time.Second)
		defer cancel()
		models, err := service.Models(ctx)
		writeResult(w, models, err)
	})
	mux.HandleFunc("/v1/models/default", func(w http.ResponseWriter, r *http.Request) {
		if !allowMutation(w, r, http.MethodPut, options) {
			return
		}
		var request struct {
			Model string `json:"model"`
		}
		if !decodeBody(w, r, &request) {
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 45*time.Second)
		defer cancel()
		err := service.SetDefaultModel(ctx, request.Model)
		writeResult(w, map[string]any{"updated": err == nil, "defaultModel": strings.TrimSpace(request.Model)}, err)
	})
	mux.HandleFunc("/v1/models/auth", func(w http.ResponseWriter, r *http.Request) {
		if !allowMutation(w, r, http.MethodPost, options) {
			return
		}
		var request struct {
			Provider string `json:"provider"`
			APIKey   string `json:"apiKey"`
		}
		if !decodeBody(w, r, &request) {
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()
		profile, err := service.SaveAPIKey(ctx, request.Provider, request.APIKey)
		writeResult(w, map[string]any{"updated": err == nil, "provider": strings.ToLower(strings.TrimSpace(request.Provider)), "profile": profile}, err)
	})
	mux.HandleFunc("/v1/agents", func(w http.ResponseWriter, r *http.Request) {
		if !allowMethod(w, r, http.MethodGet) {
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()
		agents, err := service.ListAgents(ctx)
		writeResult(w, map[string]any{"count": len(agents), "agents": agents}, err)
	})
	mux.HandleFunc("/v1/agents/sync", func(w http.ResponseWriter, r *http.Request) {
		if !allowMutation(w, r, http.MethodPost, options) {
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Minute)
		defer cancel()
		agents, err := service.SyncAgents(ctx)
		writeResult(w, map[string]any{"synced": err == nil, "count": len(agents), "agents": agents}, err)
	})
	return securityHeaders(mux)
}

func allowMethod(w http.ResponseWriter, r *http.Request, method string) bool {
	if r.Method == method {
		return true
	}
	w.Header().Set("Allow", method)
	writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "request method is not supported")
	return false
}

func allowMutation(w http.ResponseWriter, r *http.Request, method string, options HTTPOptions) bool {
	if !allowMethod(w, r, method) {
		return false
	}
	if r.Header.Get(options.MutationHeader) != options.MutationValue {
		writeError(w, http.StatusForbidden, "REQUEST_HEADER_REQUIRED", "mutation request header is missing")
		return false
	}
	return true
}

func decodeBody(w http.ResponseWriter, r *http.Request, destination any) bool {
	if !strings.HasPrefix(strings.ToLower(r.Header.Get("Content-Type")), "application/json") {
		writeError(w, http.StatusUnsupportedMediaType, "JSON_REQUIRED", "request must use application/json")
		return false
	}
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 16<<10))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "request body is invalid")
		return false
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "request must contain one JSON object")
		return false
	}
	return true
}

func writeResult(w http.ResponseWriter, value any, err error) {
	if err == nil {
		writeJSON(w, http.StatusOK, value)
		return
	}
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		writeError(w, http.StatusGatewayTimeout, "OPENCLAW_TIMEOUT", "OpenClaw operation timed out")
	case errors.Is(err, ErrUnavailable):
		writeError(w, http.StatusServiceUnavailable, "OPENCLAW_UNAVAILABLE", "OpenClaw CLI is unavailable")
	case errors.Is(err, ErrInvalidModel):
		writeError(w, http.StatusBadRequest, "INVALID_MODEL", "model ID is invalid")
	case errors.Is(err, ErrModelUnavailable):
		writeError(w, http.StatusBadRequest, "MODEL_UNAVAILABLE", "model is not available")
	case errors.Is(err, ErrInvalidProvider):
		writeError(w, http.StatusBadRequest, "INVALID_PROVIDER", "provider ID is invalid")
	case errors.Is(err, ErrInvalidAPIKey):
		writeError(w, http.StatusBadRequest, "INVALID_API_KEY", "API key is invalid")
	case errors.Is(err, ErrManifestMissing):
		writeError(w, http.StatusServiceUnavailable, "MANIFEST_UNAVAILABLE", "agent manifest is unavailable")
	default:
		writeError(w, http.StatusBadGateway, "OPENCLAW_OPERATION_FAILED", "OpenClaw operation failed")
	}
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, map[string]any{"error": map[string]string{"code": code, "message": message}})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "no-referrer")
		next.ServeHTTP(w, r)
	})
}
