package handler

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/jezrealdev/l2e-shield/internal/cache"
	"github.com/jezrealdev/l2e-shield/internal/security"
)

type ChatRequest struct {
	Contents []Content `json:"contents"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func HandleProxy(c *cache.Cache, geminiAPIKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "only POST requests are allowed"})
			return
		}

		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to read request body"})
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		var chatReq ChatRequest
		err = json.Unmarshal(bodyBytes, &chatReq)
		if err != nil || len(chatReq.Contents) == 0 || len(chatReq.Contents[0].Parts) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request payload structure"})
			return
		}

		prompt := chatReq.Contents[0].Parts[0].Text
		err = security.ValidatePrompt(prompt)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
			return
		}

		hasher := sha256.New()
		hasher.Write(bodyBytes)
		cacheKey := fmt.Sprintf("cache:response:%s", hex.EncodeToString(hasher.Sum(nil)))

		ctx := r.Context()
		cachedVal, err := c.Get(ctx, cacheKey)
		if err == nil && cachedVal != "" {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-Cache", "HIT")
			w.Write([]byte(cachedVal))
			return
		}

		geminiURL := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=%s", geminiAPIKey)
		proxyReq, err := http.NewRequestWithContext(ctx, http.MethodPost, geminiURL, bytes.NewBuffer(bodyBytes))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to create proxy request"})
			return
		}
		proxyReq.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Do(proxyReq)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "downstream API communication failure"})
			return
		}
		defer resp.Body.Close()

		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to read downstream response"})
			return
		}

		if resp.StatusCode != http.StatusOK {
			w.WriteHeader(resp.StatusCode)
			w.Write(respBytes)
			return
		}

		_ = c.Set(ctx, cacheKey, string(respBytes), 1*time.Hour)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Cache", "MISS")
		w.WriteHeader(http.StatusOK)
		w.Write(respBytes)
	}
}
