package workers

import (
	"Raven/internals/database"
	"sync"
	"testing"
)

func TestWorkerPoolDynamicExecution(t *testing.T) {
	database.InitiatDataStore()
	wp := NewWorkerPool(4, 50)

	// Test PING
	resChan1 := make(chan Response, 1)
	wp.JobQueue <- Job{Query: "PING", ResponseChan: resChan1}
	res1 := <-resChan1
	if res1.Result != "PONG\n" {
		t.Fatalf("Expected PONG\\n, got %q", res1.Result)
	}

	// Test SET with quoted string
	resChan2 := make(chan Response, 1)
	wp.JobQueue <- Job{Query: `SET config_title "Raven High Performance Engine"`, ResponseChan: resChan2}
	res2 := <-resChan2
	if res2.Result != "OK\n" {
		t.Fatalf("Expected OK\\n, got %q", res2.Result)
	}

	// Test GET quoted string
	resChan3 := make(chan Response, 1)
	wp.JobQueue <- Job{Query: `GET config_title`, ResponseChan: resChan3}
	res3 := <-resChan3
	if res3.Result != "Raven High Performance Engine\n" {
		t.Fatalf("Expected 'Raven High Performance Engine\\n', got %q", res3.Result)
	}

	// Test Nested Command Query
	resChan4 := make(chan Response, 1)
	wp.JobQueue <- Job{Query: `SET backup_title (GET config_title)`, ResponseChan: resChan4}
	res4 := <-resChan4
	if res4.Result != "OK\n" {
		t.Fatalf("Expected OK\\n, got %q", res4.Result)
	}

	resChan5 := make(chan Response, 1)
	wp.JobQueue <- Job{Query: `GET backup_title`, ResponseChan: resChan5}
	res5 := <-resChan5
	if res5.Result != "Raven High Performance Engine\n" {
		t.Fatalf("Expected 'Raven High Performance Engine\\n', got %q", res5.Result)
	}
}

func TestConcurrentWorkerDynamicJobs(t *testing.T) {
	database.InitiatDataStore()
	wp := NewWorkerPool(4, 100)

	var wg sync.WaitGroup
	numJobs := 40

	for i := 0; i < numJobs; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resChan := make(chan Response, 1)
			wp.JobQueue <- Job{
				Query:        `PING`,
				ResponseChan: resChan,
			}
			res := <-resChan
			if res.Result != "PONG\n" {
				t.Errorf("Worker returned %q, expected PONG\\n", res.Result)
			}
		}()
	}

	wg.Wait()
}
