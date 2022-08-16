package authorize

import (
	"github.com/gin-gonic/gin"
	"saas/app/constant"
	"saas/app/model"
	"saas/kernel/cache"
	"strconv"
)

func Admin(ctx *gin.Context) model.SysAdmin {

	var admin model.SysAdmin

	if Check(ctx) {
		if temp, exist := ctx.Get(constant.ContextAdmin); exist {
			admin = temp.(model.SysAdmin)
		} else {
			cache.FindById(ctx, &admin, Id(ctx))
			if admin.Id > 0 {
				ctx.Set(constant.ContextAdmin, admin)
				return admin
			}
		}
	}

	return admin
}

func Check(ctx *gin.Context) bool {
	if Id(ctx) > 0 {
		return true
	} else {
		return false
	}
}

func Id(ctx *gin.Context) int {
	var id = 0
	if ID, exist := ctx.Get(constant.ContextID); exist && ID != "" {
		if t, err := strconv.Atoi(ID.(string)); err == nil {
			id = t
		}
	}
	return id
}
