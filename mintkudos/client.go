package mintkudos

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var ErrRateLimit = errors.New("rate limited")

const (
	pageLimit = 1000
	apiHost   = "https://api.mintkudos.xyz"

	// public endpoints
	apiTokens = "/v1/tokens"

	// internal endpoints
	apiCommunities     = "/communities"
	apiCommunityTokens = "/communities/%s/tokens?isHidden=%v"
	apiMembers         = "/communities/%s/members"
	apiEndorsers       = "/token/%d/endorsers"
	apiContributors    = "/token/%d/contributors"
	apiOwners          = "/token/%d/owners"
)

type Client interface {
	GetTokens(ctx context.Context) ([]*Token, error)
	GetCommunityTokens(ctx context.Context, communityId string) (*CommunityTokens, error)
	GetCommunities(ctx context.Context) ([]*Community, error)
	GetMembers(ctx context.Context, communityID string) ([]*Member, error)
	GetEndorsers(ctx context.Context, tokenID int) ([]*Endorser, error)
	GetContributors(ctx context.Context, tokenID int) ([]*Account, error)
	GetOwners(ctx context.Context, tokenID int) ([]*Account, error)
}

type client struct{}

func (c *client) GetTokens(ctx context.Context) ([]*Token, error) {
	return getPageableList[Token](ctx, apiTokens)
}

func (c *client) GetContributors(ctx context.Context, tokenID int) ([]*Account, error) {
	return getPageableList[Account](ctx, fmt.Sprintf(apiContributors, tokenID))
}

func (c *client) GetOwners(ctx context.Context, tokenID int) ([]*Account, error) {
	return getPageableList[Account](ctx, fmt.Sprintf(apiOwners, tokenID))
}

func (c *client) GetMembers(ctx context.Context, communityID string) ([]*Member, error) {
	return getPageableList[Member](ctx, fmt.Sprintf(apiMembers, communityID))
}

func (c *client) getCommunityTokensByVisibility(ctx context.Context, communityId string, isHidden bool) ([]*CommunityToken, error) {
	return getPageableList[CommunityToken](ctx, fmt.Sprintf(apiCommunityTokens, communityId, isHidden))
}

func (c *client) GetCommunities(ctx context.Context) ([]*Community, error) {
	return getListAPI[Community](ctx, apiCommunities)
}

func (c *client) GetEndorsers(ctx context.Context, tokenID int) ([]*Endorser, error) {
	return getListAPI[Endorser](ctx, fmt.Sprintf(apiEndorsers, tokenID))
}

func unmarshalAny[T any](bytes []byte) (*T, error) {
	out := new(T)
	if err := json.Unmarshal(bytes, out); err != nil {
		return nil, err
	}
	return out, nil
}

func getListAPI[T DataType](ctx context.Context, urlPath string) ([]*T, error) {
	var result []*T
	var resp Response
	if err := get(ctx, fmt.Sprintf("%s%s", apiHost, urlPath), &resp); err != nil {
		return nil, err
	}
	for _, item := range resp.Data {
		b, err := item.MarshalJSON()
		if err != nil {
			return nil, err
		}
		res, err := unmarshalAny[T](b)
		if err != nil {
			return nil, err
		}
		result = append(result, res)
	}
	if len(result) == 0 {
		return make([]*T, 0), nil
	}
	return result, nil
}

// getPageableList iterates over all pages of the API
func getPageableList[T Pageable](ctx context.Context, urlPath string) ([]*T, error) {
	var result []*T
	var offset int64
	for {
		sym := "?"
		if strings.Contains(urlPath, "?") {
			sym = "&"
		}
		path := fmt.Sprintf("%s%soffset=%d&limit=%d", urlPath, sym, offset, pageLimit)
		items, err := getListAPI[T](ctx, path)
		if err != nil {
			return nil, err
		}

		if len(items) == 0 {
			break
		}

		result = append(result, items...)
		offset += int64(len(items))
	}
	if len(result) == 0 {
		return make([]*T, 0), nil
	}
	return result, nil
}

func (c *client) GetCommunityTokens(ctx context.Context, communityId string) (*CommunityTokens, error) {
	hidden, err := c.getCommunityTokensByVisibility(ctx, communityId, true)
	if err != nil {
		return nil, err
	}
	visible, err := c.getCommunityTokensByVisibility(ctx, communityId, false)
	if err != nil {
		return nil, err
	}
	return &CommunityTokens{CommunityID: communityId, HiddenTokens: hidden, VisibleTokens: visible}, nil
}

func NewClient() Client {
	return &client{}
}
