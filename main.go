package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"mintkudos-export/mintkudos"
)

type TokenReport struct {
	Token     *mintkudos.Token      `json:"token"`
	Endorsers []*mintkudos.Endorser `json:"endorsers"`
	Claimed   []*mintkudos.Account  `json:"claimed"`
	Issued    []*mintkudos.Account  `json:"issued"`
	IsHidden  bool                  `json:"isHidden"`
}

type CommunityReport struct {
	Community *mintkudos.Community `json:"community"`
	Members   []*mintkudos.Member  `json:"members"`
	Tokens    []*TokenReport       `json:"tokens"`
}

func writeJson(i interface{}, filename string) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	b, _ := json.Marshal(i)
	if _, err := f.WriteString(string(b)); err != nil {
		panic(err)
	}
}

func main() {
	directory := fmt.Sprintf("output-%d", time.Now().UnixNano())
	os.Mkdir(directory, os.ModePerm)
	os.Mkdir(fmt.Sprintf("%s/communities", directory), os.ModePerm)
	os.Mkdir(fmt.Sprintf("%s/tokens", directory), os.ModePerm)

	fmt.Printf("writing to %s\n", directory)

	c := mintkudos.NewClient()
	ctx := context.Background()
	tokens, err := c.GetTokens(ctx)
	if err != nil {
		log.WithError(err).Fatal("failed to load tokens")
	}

	tmux := sync.Mutex{}
	tch := make(chan *mintkudos.Token)
	grp, _ := errgroup.WithContext(ctx)
	var trs []*TokenReport
	for i := 0; i < 10; i++ {
		grp.Go(func() error {
			for t := range tch {
				logger := log.WithFields(log.Fields{
					"tokenId": t.TokenId,
				})
				e, err := c.GetEndorsers(ctx, t.TokenId)
				if err != nil {
					logger.WithError(err).Fatal("failed to load endorsers for token")
				}
				o, err := c.GetOwners(ctx, t.TokenId)
				if err != nil {
					logger.WithError(err).Fatal("failed to load owners for token")
				}
				contributors, err := c.GetContributors(ctx, t.TokenId)
				if err != nil {
					logger.WithError(err).Fatal("failed to load contributors for token")
				}
				tr := &TokenReport{
					Token:     t,
					Endorsers: e,
					Issued:    contributors,
					Claimed:   o,
				}

				if t.ClaimabilityAttributes.TotalClaimCount != nil {
					total := *t.ClaimabilityAttributes.TotalClaimCount
					if total != -1 && total != int64(len(contributors)+len(o)) {
						logger.WithError(err).Warn("token didn't add up")
					}
				}

				logger.Info("processed token")
				tmux.Lock()
				trs = append(trs, tr)
				tmux.Unlock()
			}
			return nil
		})
	}

	grp.Go(func() error {
		for _, t := range tokens {
			tch <- t
		}
		close(tch)
		return nil
	})

	if err := grp.Wait(); err != nil {
		log.WithError(err).Fatal("errgroup error")
	}

	communities, err := c.GetCommunities(ctx)
	if err != nil {
		log.WithError(err).Fatal("failed to load communities")
	}

	cch := make(chan *mintkudos.Community)
	cgrp, _ := errgroup.WithContext(ctx)
	cmux := sync.Mutex{}
	var crs []*CommunityReport
	for i := 0; i < 10; i++ {
		cgrp.Go(func() error {
			for cmty := range cch {
				logger := log.WithFields(log.Fields{
					"community": cmty.UniqId,
				})
				m, err := c.GetMembers(ctx, cmty.UniqId)
				if err != nil {
					logger.WithError(err).Fatal("failed to load members")
				}
				ct, err := c.GetCommunityTokens(ctx, cmty.UniqId)
				if err != nil {
					logger.WithError(err).Fatal("failed to load community tokens")
				}

				var ctrs []*TokenReport
				for _, trs := range trs {
					for _, ht := range ct.HiddenTokens {
						if ht.TokenId == trs.Token.TokenId {
							trs.IsHidden = true
							break
						}
					}
					if trs.Token.CommunityId == cmty.UniqId {
						ctrs = append(ctrs, trs)
					}
				}

				logger.Info("processed community")
				cmux.Lock()
				crs = append(crs, &CommunityReport{
					Community: cmty,
					Members:   m,
					Tokens:    ctrs,
				})
				cmux.Unlock()
			}
			return nil
		})
	}

	cgrp.Go(func() error {
		for _, cmty := range communities {
			cch <- cmty
		}
		close(cch)
		return nil
	})

	if err := cgrp.Wait(); err != nil {
		log.WithError(err).Fatal("community errgroup error")
	}

	for _, cr := range crs {
		writeJson(cr, fmt.Sprintf("%s/communities/%s.json", directory, cr.Community.UniqId))
	}
	for _, tr := range trs {
		writeJson(tr, fmt.Sprintf("%s/tokens/%d.json", directory, tr.Token.TokenId))
	}
}
