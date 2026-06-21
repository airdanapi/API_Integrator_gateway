package notification

import (
	"context"
	"time"
)

type AlertGenerator interface {
	GenerateAlerts(ctx context.Context) (int, error)
}

type LoggerFunc func(format string, args ...any)

func StartScheduler(ctx context.Context, generator AlertGenerator, interval time.Duration, logf LoggerFunc) {
	if generator == nil || interval <= 0 {
		return
	}
	if logf == nil {
		logf = func(string, ...any) {}
	}
	go func() {
		runGenerator(ctx, generator, logf)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				runGenerator(ctx, generator, logf)
			}
		}
	}()
}

func runGenerator(ctx context.Context, generator AlertGenerator, logf LoggerFunc) {
	created, err := generator.GenerateAlerts(ctx)
	if err != nil {
		logf("notification alert generation failed: %v", err)
		return
	}
	if created > 0 {
		logf("notification alert generation created %d alert(s)", created)
	}
}
