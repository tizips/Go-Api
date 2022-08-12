package basic

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/golang-module/carbon/v2"
	"golang.org/x/crypto/bcrypt"
	"saas/app/constant"
	"saas/app/model"
	basicForm "saas/app/request/admin/basic"
	accountResponse "saas/app/response/admin/basic"
	"saas/app/service/basic"
	helperService "saas/app/service/helper"
	"saas/kernel/authorize"
	"saas/kernel/config"
	"saas/kernel/data"
	basicResponse "saas/kernel/response"
	"strconv"
)

func DoLoginByAccount(ctx *gin.Context) {

	var request basicForm.DoLoginByAccess

	if err := ctx.ShouldBind(&request); err != nil {
		basicResponse.FailByRequest(ctx, err)
		return
	}

	var SysAdmin model.SysAdmin

	data.Database.Where("username", request.Username).Where("is_enable = ?", constant.IsEnableYes).Find(&SysAdmin)

	if SysAdmin.Id <= 0 {
		basicResponse.Fail(ctx, "用户名或密码错误")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(SysAdmin.Password), []byte(request.Password)); err != nil {
		basicResponse.Fail(ctx, "用户名或密码错误")
		return
	}

	now := carbon.Now()

	claims := jwt.RegisteredClaims{
		Issuer:    constant.ContextAdmin,
		Subject:   strconv.Itoa(SysAdmin.Id),
		NotBefore: jwt.NewNumericDate(now.Carbon2Time()),
		IssuedAt:  jwt.NewNumericDate(now.Carbon2Time()),
		ExpiresAt: jwt.NewNumericDate(now.AddHours(config.Values.Jwt.Lifetime).Carbon2Time()),
		ID:        helperService.JwtToken(SysAdmin.Id),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(config.Values.Jwt.Secret))

	if err != nil {
		basicResponse.Fail(ctx, "用户名或密码错误")
		return
	}

	basicResponse.SuccessByData(ctx, accountResponse.DoLoginByAccess{
		Token:    signed,
		ExpireAt: now.AddHours(config.Values.Jwt.Lifetime).Timestamp(),
	})

}

func DoLoginByQrcode(ctx *gin.Context) {

}

func DoLogout(ctx *gin.Context) {

	claims := authorize.Jwt(ctx)

	ok := basic.BlackJwt(ctx, constant.ContextAdmin, *claims)

	if !ok {
		basicResponse.Fail(ctx, "退出失败，请稍后重试！")
		return
	}

	basicResponse.Success(ctx)
}
