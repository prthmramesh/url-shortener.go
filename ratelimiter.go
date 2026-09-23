package main

import (
	"context"
	"log/slog"
	"time"
)

func allowRequest(shortCode string) bool {
	rateLimitMutex.Lock()
	defer rateLimitMutex.Unlock()

	cutoff := time.Now().Add(-time.Minute)
	valid := requestLog[shortCode][:0]

	for _, t := range requestLog[shortCode] {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}

	if len(valid) >= 5 {
		requestLog[shortCode] = valid
		return false
	}

	requestLog[shortCode] = append(valid, time.Now())
	return true
}

func requestLogCleaner(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			cutoff := time.Now().Add(-time.Minute)

			rateLimitMutex.Lock()
			for shortUrl := range requestLog {
				valid := requestLog[shortUrl][:0]
				for _, t := range requestLog[shortUrl] {
					if t.After(cutoff) {
						valid = append(valid, t)
					}
				}

				if len(valid) == 0 {
					delete(requestLog, shortUrl)
					slog.Info("request_log entry deleted", "short_code", shortUrl)
				} else {
					requestLog[shortUrl] = valid
				}
			}
			rateLimitMutex.Unlock()

		case <-ctx.Done():
			slog.Warn("main context cancelled")
			return
		}
	}
}
