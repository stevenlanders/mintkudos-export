package mintkudos

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cenkalti/backoff/v4"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func withBackoff(
	ctx context.Context, operation func(ctx context.Context) error,
) error {
	bo := backoff.NewExponentialBackOff()
	bo.InitialInterval = 5 * time.Second
	bo.MaxInterval = 30 * time.Second

	err := backoff.Retry(func() error {
		if ctx.Err() != nil {
			return backoff.Permanent(ctx.Err())
		}
		err := operation(ctx)
		if err == nil {
			return nil
		} else if ctx.Err() != nil {
			return backoff.Permanent(ctx.Err())
		} else if err == ErrRateLimit {
			log.Warnf("%s (retryable)", err.Error())
			return err
		}
		return backoff.Permanent(err)
	}, bo)
	if err != nil {
		log.Errorf("failed with error: %s", err.Error())
	}
	return err
}

func get(ctx context.Context, path string, target interface{}) error {
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return err
	}

	return withBackoff(ctx, func(ctx context.Context) error {
		req = req.WithContext(ctx)

		dc := http.DefaultClient
		resp, err := dc.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode == 429 || resp.StatusCode == 502 {
			return ErrRateLimit
		}
		if resp.StatusCode >= 400 {
			return fmt.Errorf("error code %d, path=%s", resp.StatusCode, path)
		}
		return json.NewDecoder(resp.Body).Decode(target)
	})
}
