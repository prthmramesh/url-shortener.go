package main

import (
	"context"
	"log/slog"
)

func clickRecorder(ctx context.Context, clicks <-chan ClickEvent) {
	for {
		select {
		case <-ctx.Done():
			slog.Info("context closed", "reason", "shutdown")
			return
		case click, ok := <-clicks:
			if !ok {
				slog.Warn("click_channel_closed")
				return
			}

			clickCounterMutex.Lock()
			clickCounter[click.ShortCode]++
			count := clickCounter[click.ShortCode]
			clickCounterMutex.Unlock()

			slog.Info("short_code called", "current_count", count, "click_code", click.ShortCode)
		}
	}
}
