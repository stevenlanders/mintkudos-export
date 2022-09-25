package mintkudos

import (
	"encoding/json"
	"errors"
)

// These are generated from example json

type Attribute struct {
	FieldName string `json:"fieldName"`
	Type      string `json:"type"`
	Value     string `json:"value"`
}

type ClaimabilityAttributes struct {
	IsSignatureRequired bool   `json:"isSignatureRequired"`
	IsAllowlistRequired bool   `json:"isAllowlistRequired"`
	TotalClaimCount     *int64 `json:"totalClaimCount"`
	RemainingClaimCount *int64 `json:"remainingClaimCount"`
	ExpirationTimestamp *int64 `json:"expirationTimestamp"`
}

type Token struct {
	TokenId                int                     `json:"tokenId"`
	Headline               string                  `json:"headline"`
	Description            string                  `json:"description"`
	StartDateTimestamp     *int64                  `json:"startDateTimestamp"`
	EndDateTimestamp       *int64                  `json:"endDateTimestamp"`
	Links                  []string                `json:"links"`
	CommunityId            string                  `json:"communityId"`
	CreatedByAddress       string                  `json:"createdByAddress"`
	CreatedAtTimestamp     int64                   `json:"createdAtTimestamp"`
	ImageUrl               string                  `json:"imageUrl"`
	ClaimabilityAttributes *ClaimabilityAttributes `json:"claimabilityAttributes"`
	CustomAttributes       []*Attribute            `json:"customAttributes"`
}

type Endorser struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

type CommunityToken struct {
	TokenId  int    `json:"tokenId"`
	Headline string `json:"headline"`
	ImageUrl string `json:"imageUrl"`
}

type CommunityTokens struct {
	CommunityID   string            `json:"communityId"`
	VisibleTokens []*CommunityToken `json:"visibleTokens"`
	HiddenTokens  []*CommunityToken `json:"hiddenTokens"`
}

type Community struct {
	UniqId      string `json:"uniqId"`
	DisplayId   string `json:"displayId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	LogoUrl     string `json:"logoUrl"`
}

// Account represents an eth account that can receive or claim token
type Account struct {
	Id            string `json:"id"`
	Username      string `json:"username"`
	Role          string `json:"role"`
	AccountUniqId string `json:"accountUniqId"`
	AccountType   string `json:"accountType"`
}

func (ct *CommunityTokens) IsHidden(tokenID int) (bool, error) {
	for _, t := range ct.HiddenTokens {
		if t.TokenId == tokenID {
			return true, nil
		}
	}
	for _, t := range ct.VisibleTokens {
		if t.TokenId == tokenID {
			return false, nil
		}
	}
	return false, errors.New("not found in community list")
}

type Role struct {
	Role string `json:"role"`
}

type Member struct {
	Username       string  `json:"username"`
	Id             string  `json:"id"`
	CommunityRoles []*Role `json:"communityRoles"`
}

//Pageable types have limit/offset functionality and can be iterated
type Pageable interface {
	CommunityToken | Token | Member | Account
}

//DataType types are represented in the .Data[] fields of responses
type DataType interface {
	Pageable | Endorser | Community
}

type Response struct {
	Data []*json.RawMessage `json:"Data"`
}
