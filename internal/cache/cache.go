package cache

import (
	"sync"
	"time"

	models "github.com/AngelPwG/devprofile/internal/domain"
)

const TTL = time.Hour

type Cache struct {
	mu    sync.RWMutex
	store map[string]*models.Profile
}

var Instance = &Cache{
	store: make(map[string]*models.Profile),
}

func (c *Cache) Get(githubUser string) (*models.Profile, bool) {
	c.mu.Lock()
	defer c.mu.RUnlock()

	profile, exists := c.store[githubUser]
	if !exists {
		return nil, false
	}

	profileLastUpdated, _ := time.Parse(time.RFC3339, profile.UpdatedAt)

	if time.Since(profileLastUpdated) > TTL {
		return profile, false
	}

	return profile, true
}

func (c *Cache) Set(profile *models.Profile) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[profile.GithubUser] = profile
}

func (c *Cache) MinutesUntilRefresh(githubUser string) int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	profile, exists := c.store[githubUser]
	if !exists {
		return 0
	}

	updatedAt, err := time.Parse(time.RFC3339, profile.UpdatedAt)
	if err != nil {
		return 0
	}
	remaining := TTL - time.Since(updatedAt)
	if remaining <= 0 {
		return 0
	}

	return int(remaining.Minutes()) + 1
}
