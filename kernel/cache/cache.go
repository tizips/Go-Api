package cache

import (
	"fmt"
	"saas/kernel/app"
	"time"
)

func Key(table string, id any) string {
	return fmt.Sprintf("%s:%s:%s:%v", app.Cfg.Server.Name, "cache", table, id)
}

func ttl() time.Duration {
	return time.Duration(86400) * time.Second
}
