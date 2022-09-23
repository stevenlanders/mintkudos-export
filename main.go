package main

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/sync/errgroup"
	"mintkudos-export/mintkudos"
	"os"
	"time"
)

func writeJson(i interface{}, filename string) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	b, _ := json.Marshal(i)
	if _, err := f.WriteString(string(b)); err != nil {
		panic(err)
	}
}

func main() {
	c := mintkudos.NewClient()
	ctx := context.Background()
	tokens, err := c.GetTokens(ctx)
	if err != nil {
		panic(err)
	}

	ch := make(chan *mintkudos.Token)
	grp, ctx := errgroup.WithContext(ctx)

	directory := fmt.Sprintf("output-%d", time.Now().UnixNano())

	os.Mkdir(directory, os.ModePerm)
	os.Mkdir(fmt.Sprintf("%s/owners", directory), os.ModePerm)

	fmt.Printf("writing to %s", directory)

	for i := 0; i < 10; i++ {
		grp.Go(func() error {
			for t := range ch {
				owners, err := c.GetOwners(ctx, t.TokenId)
				if err != nil {
					panic(err)
				}
				if owners == nil {
					owners = []*mintkudos.Owner{}
				}
				writeJson(owners, fmt.Sprintf("%s/owners/owners-%d.json", directory, t.TokenId))
			}
			return nil
		})
	}

	for _, t := range tokens {
		ch <- t
	}
	close(ch)

	writeJson(tokens, fmt.Sprintf("%s/tokens.json", directory))
}
