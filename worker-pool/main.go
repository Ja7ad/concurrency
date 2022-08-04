package main

import (
	"context"
	"fmt"
	"log"
	"sync"
)

type Job func(ctx context.Context, wg *sync.WaitGroup, name string) result

type result struct {
	message string
	err     error
}

func main() {
	const numbJobs = 8
	jobs := make(chan Job, numbJobs)
	results := make(chan result, numbJobs)
	names := []string{"A", "B", "C", "D", "E", "F", "G", "H"}
	ctx := context.Background()

	for i, name := range names {
		go workerEfficient(ctx, i, name, jobs, results)
	}

	for i := 0; i < numbJobs; i++ {
		jobs <- greetingJob
	}

	close(jobs)
	for a := 1; a <= numbJobs; a++ {
		fmt.Println(<-results)
	}
	close(results)
}

func workerEfficient(ctx context.Context, id int, name string, jobs <-chan Job, results chan<- result) {
	wg := &sync.WaitGroup{}
	for j := range jobs {
		wg.Add(1)
		log.Printf("job %d started do job", id)
		go func(ctx context.Context, wg *sync.WaitGroup, job Job, name string) {
			results <- job(ctx, wg, name)
		}(ctx, wg, j, name)
	}
	wg.Wait()
}

func greetingJob(ctx context.Context, wg *sync.WaitGroup, name string) result {
	defer wg.Done()
	select {
	case <-ctx.Done():
		return result{"", ctx.Err()}
	default:
	}
	return result{fmt.Sprintf("greeting %s", name), nil}
}
