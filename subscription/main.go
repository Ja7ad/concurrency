package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
)

const exampleAPIAddress = "https://random-data-api.com/api/stripe/random_stripe"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	sub := NewSubscription(ctx, NewFetcher(exampleAPIAddress), 3)

	time.AfterFunc(1*time.Minute, func() {
		cancel()
		log.Println("canceled subscription task")
		os.Exit(0)
	})

	for card := range sub.Updates() {
		fmt.Println(card)
	}
}
