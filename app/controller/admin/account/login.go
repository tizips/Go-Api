package account

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon/v2"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"saas/app/constant"
	"saas/app/form/admin/account"
	"saas/app/model"
	accountResponse "saas/app/response/admin/account"
	helperService "saas/app/service/helper"
	"saas/kernel/config"
	"saas/kernel/data"
	"saas/kernel/response"
)

func DoLoginByAccount(ctx *gin.Context) {

	var params account.DoLoginForm

	if err := ctx.ShouldBind(&params); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: err.Error(),
		})
		return
	}

	var SysAdmin model.SysAdmin

	data.Database.Where("username", params.Username).Where("is_enable = ?", constant.IsEnableYes).First(&SysAdmin)

	if SysAdmin.Id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40100,
			Message: "用户名或密码错误",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(SysAdmin.Password), []byte(params.Password)); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40100,
			Message: "用户名或密码错误",
		})
		return
	}

	now := carbon.Now()

	claims := jwt.StandardClaims{
		Issuer:    "admin",
		Id:        fmt.Sprintf("%d", SysAdmin.Id),
		NotBefore: now.Timestamp(),
		IssuedAt:  now.Timestamp(),
		ExpiresAt: now.AddHours(12).Timestamp(),
		Audience:  helperService.JwtToken(SysAdmin.Id),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(config.Configs.Jwt.Secret))

	if err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    50000,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Code:    20000,
		Message: "Success",
		Data: accountResponse.LoginResponse{
			Token:    signed,
			ExpireAt: now.AddHours(12).Timestamp(),
		},
	})
}

func DoLoginByQrcode(ctx *gin.Context) {

}
