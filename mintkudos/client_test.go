package mintkudos

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetCommunityTokens(t *testing.T) {
	c := NewClient()
	ctx := context.Background()
	ct, err := c.GetCommunityTokens(ctx, "vectorDAO")
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(ct.HiddenTokens), 0)
	assert.GreaterOrEqual(t, len(ct.VisibleTokens), 4)
}

func TestClient_GetOwners(t *testing.T) {
	c := NewClient()
	ctx := context.Background()
	owners, err := c.GetOwners(ctx, 1577)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(owners), 10)
}
