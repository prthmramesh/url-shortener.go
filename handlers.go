package main

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/lib/pq"
)

func urlshortener(w http.ResponseWriter, r *http.Request) {
	var urlCode ShortURL

	if err := json.NewDecoder(r.Body).Decode(&urlCode); err != nil {
		http.Error(w, "url shortener request rejected", http.StatusBadRequest)
		slog.Warn("request rejected", "reason", "error while decoding json", "error", err.Error())
		return
	}

	if urlCode.OriginalURL == "" {
		http.Error(w, "Url field cannot be empty", http.StatusBadRequest)
		slog.Warn("request rejected", "reason", "url field empty", "value", urlCode.OriginalURL)
		return
	}

	if _, err := url.ParseRequestURI(urlCode.OriginalURL); err != nil {
		http.Error(w, "invalid url sent", http.StatusBadRequest)
		slog.Warn("request rejected", "reason", "invalid url sent", "value", urlCode.OriginalURL)
		return
	}

	shorturl := generateCode()

	urlCode.ShortCode = shorturl
	urlCode.CreatedTime = time.Now()

	if err := insertURL(urlCode.ShortCode, urlCode.OriginalURL); err != nil {

		if pq, ok := err.(*pq.Error); ok && pq.Code == "23505" {
			http.Error(w, "short url already exists", http.StatusConflict)
			slog.Warn("db insert failed", "reason", "short url already exists", "shortened_url", urlCode.ShortCode, "original_url", urlCode.OriginalURL)
			return
		} else {
			http.Error(w, "failed to save url", http.StatusInternalServerError)
			slog.Error("db insert failed", "error", err.Error())
		}

		return
	}

	slog.Info("url shortened", "shortened_url", urlCode.ShortCode, "original_url", urlCode.OriginalURL)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(urlCode); err != nil {
		slog.Error("failed to encode response", "error", err.Error())
	}
}

func generateCode() string {
	charset := "abcdefghijklmnopqrstuvwxyz"
	random_character := make([]byte, 0, 6)

	for i := 0; i < 6; i++ {
		random_character = append(random_character, charset[rand.Intn(len(charset))])
	}

	return string(random_character)
}

func redirectUrl(w http.ResponseWriter, r *http.Request) {
	shortCode := r.PathValue("shortCode")

	entry, err := getOriginalURL(shortCode)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Short url not found", http.StatusNotFound)
			slog.Warn("short_url not found", "short_code", shortCode)
			return
		} else {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			slog.Error("db query failed", "short_code", shortCode, "error", err.Error())
		}
		return

	}

	slog.Info("redirect served", "short_code", shortCode, "destination", entry)

	eventClick <- ClickEvent{ShortCode: shortCode, Timestamp: time.Now()}

	if !allowRequest(shortCode) {
		http.Error(w, "too many requests", http.StatusTooManyRequests)
		slog.Warn("fetch call exceeded", "short_code", shortCode)
		return
	}

	http.Redirect(w, r, entry, http.StatusFound)
}

func clickEventChecker(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("shortCode")

	clickCounterMutex.RLock()
	counter, ok := clickCounter[code]
	clickCounterMutex.RUnlock()

	if !ok {
		http.Error(w, "short code not found", http.StatusNotFound)
		slog.Info("short_code_not_found", "value", code)
		return
	}

	slog.Info("short code counter fetched", "count", counter)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"shortCode": code,
		"clicks":    counter,
	})
}
