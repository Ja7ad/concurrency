package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Subscription interface {
	Updates() <-chan Card
}

type Fetcher interface {
	Fetch() (Card, error)
}

type sub struct {
	fetcher Fetcher
	updates chan Card
}

type fetcher struct {
	url string
}

type fetchResult struct {
	fetchedCard Card
	err         error
}

// NewSubscription create subscription for fetch data per freq time in second
func NewSubscription(ctx context.Context, fetcher Fetcher, freq uint) Subscription {
	s := &sub{
		fetcher: fetcher,
		updates: make(chan Card),
	}
	go s.serve(ctx, freq)
	return s
}

func NewFetcher(url string) Fetcher {
	return &fetcher{
		url: url,
	}
}

func (f *fetcher) Fetch() (Card, error) {
	return requestAPI(f.url)
}

func (s *sub) serve(ctx context.Context, freq uint) {
	ticker := time.NewTicker(time.Duration(freq) * time.Second)
	done := make(chan fetchResult, 1)

	var (
		fetchedCard         Card
		fetchResponseStream chan Card
		pending             bool
	)

	for {

		if pending {
			fetchResponseStream = s.updates
		} else {
			fetchResponseStream = nil
		}

		select {
		case <-ticker.C:
			if pending {
				break
			}
			go func() {
				fetched, err := s.fetcher.Fetch()
				done <- fetchResult{fetched, err}
			}()
		case result := <-done:
			fetchedCard = result.fetchedCard
			if result.err != nil {
				log.Printf("fetch got error %v", result.err)
				break
			}
			pending = true
		case fetchResponseStream <- fetchedCard:
			pending = false
		case <-ctx.Done():
			return
		}
	}
}

func (s *sub) Updates() <-chan Card {
	return s.updates
}

func requestAPI(url string) (Card, error) {
	card := Card{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return Card{}, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return Card{}, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return Card{}, err
	}
	if err := json.Unmarshal(body, &card); err != nil {
		return Card{}, err
	}
	return card, nil
}
