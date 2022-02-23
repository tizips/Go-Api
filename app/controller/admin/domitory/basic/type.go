package basic

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"saas/app/form/admin/dormitory/basic"
	"saas/app/model"
	basicResponse "saas/app/response/admin/dormitory/basic"
	"saas/kernel/data"
	"saas/kernel/response"
	"strconv"
)

func DoTypeByCreate(ctx *gin.Context) {

	var former basic.DoTypeByCreateForm
	if err := ctx.ShouldBind(&former); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: err.Error(),
		})
		return
	}

	typ := model.DorType{
		Name:     former.Name,
		Order:    former.Order,
		IsEnable: former.IsEnable,
	}

	if data.Database.Create(&typ); typ.Id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: "添加失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Code:    20000,
		Message: "Success",
	})

}

func DoTypeByUpdate(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: "楼栋ID获取失败",
		})
		return
	}

	var former basic.DoTypeByUpdateForm
	if err := ctx.ShouldBind(&former); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: err.Error(),
		})
		return
	}

	var typ model.DorType
	data.Database.First(&typ, id)
	if typ.Id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "未找到该楼栋",
		})
		return
	}

	typ.Name = former.Name
	typ.Order = former.Order
	typ.IsEnable = former.IsEnable

	if t := data.Database.Save(&typ); t.RowsAffected <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: "修改失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Code:    20000,
		Message: "Success",
	})

}

func DoTypeByDelete(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: "楼栋ID获取失败",
		})
		return
	}

	var typ model.DorType
	data.Database.First(&typ, id)
	if typ.Id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "未找到该楼栋",
		})
		return
	}

	if t := data.Database.Delete(&typ); t.RowsAffected <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: "删除失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Code:    20000,
		Message: "Success",
	})

}

func DoTypeByEnable(ctx *gin.Context) {

	var former basic.DoTypeByEnableForm
	if err := ctx.ShouldBind(&former); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: err.Error(),
		})
		return
	}

	var typ model.DorType
	data.Database.First(&typ, former.Id)
	if typ.Id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "未找到该楼栋",
		})
		return
	}

	typ.IsEnable = former.IsEnable

	if t := data.Database.Save(&typ); t.RowsAffected <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: "启禁失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Code:    20000,
		Message: "Success",
	})

}

func ToTypeByList(ctx *gin.Context) {

	responses := response.Responses{
		Code:    20000,
		Message: "Success",
		Data:    []interface{}{},
	}

	var types []model.DorType
	data.Database.Order("`order` asc").Order("`id` desc").Find(&types)

	for _, item := range types {
		responses.Data = append(responses.Data, basicResponse.ToTypeByListResponse{
			Id:        item.Id,
			Name:      item.Name,
			Order:     item.Order,
			IsEnable:  item.IsEnable,
			CreatedAt: item.CreatedAt.ToDateTimeString(),
		})
	}

	ctx.JSON(http.StatusOK, responses)
}

func ToTypeByOnline(ctx *gin.Context) {

	responses := response.Responses{
		Code:    20000,
		Message: "Success",
		Data:    []interface{}{},
	}

	var types []model.DorType
	data.Database.Order("`order` asc").Order("`id` desc").Find(&types)

	for _, item := range types {
		responses.Data = append(responses.Data, basicResponse.ToTypeByOnlineResponse{
			Id:   item.Id,
			Name: item.Name,
		})
	}

	ctx.JSON(http.StatusOK, responses)
}
