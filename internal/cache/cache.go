package cache

import "time"

const TTL = time.Hour

func CanRefresh(UpdatedAt string) (bool, int, error) {
	t, err := time.Parse(time.RFC3339, UpdatedAt)
	if err != nil {
		return false, 0, err
	}
	remaining := TTL - time.Since(t)
	if remaining <= 0 {
		return true, 0, nil
	}
	return false, int(remaining.Seconds()), nil
}
