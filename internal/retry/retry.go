package retry

import (
	"context"
	"fmt"
	"time"
)

func Do(ctx context.Context, attempts int, baseDelay time.Duration, fn func() error) error {
	for i := attempts; i >= 0; i-- {
		err := fn()
		if err == nil {
			return nil
		}
		ctx.Done()
		contErr := ctx.Err()
		if contErr != nil {
			return fmt.Errorf("Context error: %w", err)
		}
		time.Sleep(baseDelay)
		baseDelay = 2 * baseDelay
		if i == 0 {
			return err
		}
	}
	return nil
}
