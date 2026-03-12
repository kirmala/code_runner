package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewClusterClient(addrs []string, password string) (*redis.ClusterClient, error) {
    rdb := redis.NewClusterClient(&redis.ClusterOptions{
        Addrs:    addrs,    
        Password: password,
    })
    var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("pinging redis: %w", err)
	}
    return rdb, nil
}