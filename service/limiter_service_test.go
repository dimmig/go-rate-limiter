package service

import (
	"go-rate-limiter/repository"
	"sync"
	"sync/atomic"
	"testing"
)

func TestConcurrentRequests(t *testing.T) {
	repo := repository.NewMemoryRepo()
	svc := NewLimiterService(repo)

	limit := 50
	totalRequests := 100

	var wg sync.WaitGroup
	var allowed int64
	var denied int64

	for i := 0; i < totalRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ok, _, err := svc.Allow("user_123", limit, 60)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if ok {
				atomic.AddInt64(&allowed, 1)
			} else {
				atomic.AddInt64(&denied, 1)
			}
		}()
	}
	wg.Wait()

	allowedCount := atomic.LoadInt64(&allowed)
	deniedCount := atomic.LoadInt64(&denied)

	if allowedCount != int64(limit) {
		t.Fatalf("expected %d allowed, got %d", limit, allowedCount)
	}
	if deniedCount != int64(totalRequests-limit) {
		t.Fatalf("expected %d denied, got %d", totalRequests-limit, deniedCount)
	}
}

func TestAllowExceedsLimit(t *testing.T) {
	repo := repository.NewMemoryRepo()
	svc := NewLimiterService(repo)

	// First 3 requests allowed
	for i := 0; i < 3; i++ {
		allowed, _, err := svc.Allow("user_1", 3, 60)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !allowed {
			t.Fatalf("request %d should be allowed", i+1)
		}
	}

	// 4th request denied
	allowed, _, err := svc.Allow("user_1", 3, 60)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if allowed {
		t.Fatal("4th request should be denied")
	}
}

func TestRemainingCalculation(t *testing.T) {
	repo := repository.NewMemoryRepo()
	svc := NewLimiterService(repo)

	_, remaining, _ := svc.Allow("user_x", 10, 60)
	if remaining != 9 {
		t.Fatalf("expected remaining=9, got %d", remaining)
	}

	for i := 0; i < 8; i++ {
		svc.Allow("user_x", 10, 60)
	}

	_, remaining, _ = svc.Allow("user_x", 10, 60)
	if remaining != 0 {
		t.Fatalf("expected remaining=0, got %d", remaining)
	}
}
