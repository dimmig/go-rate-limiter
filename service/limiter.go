package service

import (
	"context"
	"go-rate-limiter/repository"
	"time"
)

type LimiterService struct {
	repo    repository.RepoInterface
	timeout time.Duration
}

func NewLimiterService(repo repository.RepoInterface) *LimiterService {
	return &LimiterService{repo: repo, timeout: 50 * time.Millisecond}
}

func (l *LimiterService) Allow(key string, limit, expiry int) (bool, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), l.timeout)
	defer cancel()

	currentCount, err := l.repo.Increment(ctx, key, expiry)
	if err != nil {
		return false, 0, err
	}

	remaining := limit - currentCount
	if remaining < 0 {
		remaining = 0
	}

	return currentCount <= limit, remaining, nil
}
