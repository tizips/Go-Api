package cache

import (
	"fmt"
	"saas/kernel/config"
	"time"
)

func Key(table string, id any) string {
	return fmt.Sprintf("%s:%s:%s:%v", config.Values.Server.Name, config.Values.Redis.CachePrefix, table, id)
}

func ttl() time.Duration {
	return time.Duration(config.Values.Redis.CacheTtl) * time.Second
}
