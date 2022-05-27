package basic

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
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

	claims := jwt.StandardClaims{
		Issuer:    "admin",
		Id:        fmt.Sprintf("%d", SysAdmin.Id),
		NotBefore: now.Timestamp(),
		IssuedAt:  now.Timestamp(),
		ExpiresAt: now.AddHours(config.Values.Jwt.Lifetime).Timestamp(),
		Audience:  helperService.JwtToken(SysAdmin.Id),
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

	ok := basic.BlackJwt(ctx, "admin", *claims)

	if !ok {
		basicResponse.Fail(ctx, "退出失败，请稍后重试！")
		return
	}

	basicResponse.Success(ctx)
}
