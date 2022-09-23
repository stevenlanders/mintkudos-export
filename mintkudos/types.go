package mintkudos

import "time"

// These are generated from example json

type Token struct {
	TokenId                int      `json:"tokenId"`
	Headline               string   `json:"headline"`
	Description            string   `json:"description"`
	StartDateTimestamp     *int64   `json:"startDateTimestamp"`
	EndDateTimestamp       *int64   `json:"endDateTimestamp"`
	Links                  []string `json:"links"`
	CommunityId            string   `json:"communityId"`
	CreatedByAddress       string   `json:"createdByAddress"`
	CreatedAtTimestamp     int64    `json:"createdAtTimestamp"`
	ImageUrl               string   `json:"imageUrl"`
	ClaimabilityAttributes struct {
		IsSignatureRequired bool   `json:"isSignatureRequired"`
		IsAllowlistRequired bool   `json:"isAllowlistRequired"`
		TotalClaimCount     *int64 `json:"totalClaimCount"`
		RemainingClaimCount *int64 `json:"remainingClaimCount"`
		ExpirationTimestamp *int64 `json:"expirationTimestamp"`
	} `json:"claimabilityAttributes"`
	CustomAttributes []interface{} `json:"customAttributes"`
}

type TokenResponse struct {
	Limit  int      `json:"limit"`
	Offset int      `json:"offset"`
	Data   []*Token `json:"data"`
}

type Owner struct {
	WalletAddress string    `json:"walletAddress"`
	MintedAt      time.Time `json:"mintedAt"`
}

type OwnerResponse struct {
	Limit  int      `json:"limit"`
	Offset int      `json:"offset"`
	Data   []*Owner `json:"data"`
}
