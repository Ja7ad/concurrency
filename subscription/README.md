# Subscription Pattern

## abstract
This pattern is based on the Advanced Go Concurrency Patterns Talk presented at Google I/O 2013.

## Pattern Usages
1. Per frequency time, consume data from publisher pub/sub
2. Fetch data from API

## Composition patterns
- [For-Select-Done](../for-select-done)

## analyze of code
1. For getting data, we abstract Subscription interface with Updates method 

```go
type Subscription interface {
	Updates() <-chan Card
}
```
2. Then abstract Fetcher interface to get data from API with Fetch method

```go
type Fetcher interface {
	Fetch() (Card, error)
}
```

3. Now create NewSubscription and pass context, NewFetcher and freq time for fetch

```go
	ctx, cancel := context.WithCancel(context.Background())
	sub := NewSubscription(ctx, NewFetcher(exampleAPIAddress), 3)
```

4. After initiate sub struct we run serve in goroutine.

```go
// NewSubscription create subscription for fetch data per freq time in second
func NewSubscription(ctx context.Context, fetcher Fetcher, freq uint) Subscription {
	s := &sub{
		fetcher: fetcher,
		updates: make(chan Card),
	}
	go s.serve(ctx, freq)
	return s
}
```
5. To understand what's going on, let's break it down:

   - The first select case will be activated regularly by the time ticker. The fetcher will be run inside another goroutine, so that we won't be blocked if it takes a long time.
   - When the fetch result is finally ready, it will be received at the second select case. When there is an error, it breaks from the select statement and waits for the next iteration.
   - The third case is the standard case that waits for a context to terminate.