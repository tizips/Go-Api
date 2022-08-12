package snowflake

import (
	"github.com/bwmarrin/snowflake"
	"saas/kernel/config"
)

var Snowflake *snowflake.Node

func InitSnowflake() (err error) {

	Snowflake, err = snowflake.NewNode(config.Values.Server.Node)

	return
}
