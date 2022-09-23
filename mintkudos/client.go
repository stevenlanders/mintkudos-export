package mintkudos

import (
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
)

var ErrRateLimit = errors.New("rate limited")

const (
	pageLimit     = 1000
	apiHost       = "https://api.mintkudos.xyz"
	apiTokensList = "/v1/tokens?limit=%d&offset=%d"
	apiOwnersList = "/v1/tokens/%d/owners?limit=%d&offset=%d"
)

type Client interface {
	GetTokens(ctx context.Context) ([]*Token, error)
	GetOwners(ctx context.Context, tokenID int) ([]*Owner, error)
}

type client struct{}

func (c *client) GetTokens(ctx context.Context) ([]*Token, error) {
	var result []*Token
	var offset int64

	for {
		var resp TokenResponse
		path := fmt.Sprintf(apiTokensList, pageLimit, offset)
		if err := get(ctx, fmt.Sprintf("%s%s", apiHost, path), &resp); err != nil {
			return nil, err
		}
		if len(resp.Data) == 0 {
			log.WithFields(log.Fields{})
			break
		}
		result = append(result, resp.Data...)
		offset += int64(len(resp.Data))
	}
	return result, nil
}

func (c *client) GetOwners(ctx context.Context, tokenID int) ([]*Owner, error) {
	var result []*Owner
	var offset int64

	for {
		var resp OwnerResponse
		path := fmt.Sprintf(apiOwnersList, tokenID, pageLimit, offset)
		if err := get(ctx, fmt.Sprintf("%s%s", apiHost, path), &resp); err != nil {
			return nil, err
		}
		if len(resp.Data) == 0 {
			log.WithFields(log.Fields{})
			break
		}
		result = append(result, resp.Data...)
		offset += int64(len(resp.Data))
	}
	return result, nil
}

func NewClient() Client {
	return &client{}
}
