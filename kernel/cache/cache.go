package cache

import (
	"fmt"
	"saas/kernel/config"
	"time"
)

func Key(table string, id any) string {
	return fmt.Sprintf("%s:%s:%v", config.Values.Cache.Prefix, table, id)
}

func ttl() time.Duration {
	return time.Duration(config.Values.Cache.Ttl) * time.Second
}
