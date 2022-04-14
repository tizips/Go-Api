package auth

import (
	"github.com/gin-gonic/gin"
	"saas/app/constant"
	"saas/app/model"
	"saas/kernel/cache"
	"strconv"
)

func Admin(ctx *gin.Context) model.SysAdmin {
	if Check(ctx) {
		if admin, exist := ctx.Get(constant.ContextAdmin); exist {
			return admin.(model.SysAdmin)
		} else {
			var SysAdmin model.SysAdmin
			cache.Find(ctx, Id(ctx), &SysAdmin)
			if SysAdmin.Id > 0 {
				ctx.Set(constant.ContextAdmin, SysAdmin)
				return SysAdmin
			}
		}
	}
	return model.SysAdmin{}
}

func Check(ctx *gin.Context) bool {
	if Id(ctx) > 0 {
		return true
	} else {
		return false
	}
}

func Id(ctx *gin.Context) uint {
	var id uint = 0
	if ID, exist := ctx.Get(constant.ContextID); exist && ID != "" {
		if t, err := strconv.Atoi(ID.(string)); err == nil {
			id = uint(t)
		}
	}
	return id
}
