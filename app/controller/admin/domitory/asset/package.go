package asset

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"saas/app/form/admin/dormitory/asset"
	"saas/app/model"
	assetResponse "saas/app/response/admin/dormitory/asset"
	"saas/kernel/data"
	"saas/kernel/response"
	"strconv"
)

func DoPackageByCreate(ctx *gin.Context) {

	var former asset.DoPackageByCreateForm
	if err := ctx.ShouldBind(&former); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: err.Error(),
		})
		return
	}

	var deviceIds = make([]uint, len(former.Devices))

	for idx, item := range former.Devices {
		deviceIds[idx] = item.Device
	}

	var devices []model.DorDevice
	data.Database.Where("id in (?)", deviceIds).Find(&devices)

	for _, item := range former.Devices {
		mark := true
		for _, value := range devices {
			if item.Device == value.Id {
				mark = false
			}
		}
		if mark {
			ctx.JSON(http.StatusOK, response.Response{
				Code:    40400,
				Message: "部分设备未找到",
			})
			return
		}
	}

	if len(devices) != len(former.Devices) {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "部分设备选择重复",
		})
		return
	}

	tx := data.Database.Begin()

	pack := model.DorPackage{
		Name: former.Name,
	}

	if t := tx.Create(&pack); t.RowsAffected <= 0 {
		tx.Rollback()
		ctx.JSON(http.StatusOK, response.Response{
			Code:    60000,
			Message: "添加失败",
		})
		return
	}

	var bindings []model.DorPackageDetail
	for _, item := range former.Devices {
		bindings = append(bindings, model.DorPackageDetail{
			PackageId: pack.Id,
			DeviceId:  item.Device,
			Number:    item.Number,
		})
	}

	if t := tx.Create(&bindings); t.RowsAffected <= 0 {
		tx.Rollback()
		ctx.JSON(http.StatusOK, response.Response{
			Code:    60000,
			Message: "添加失败",
		})
		return
	}

	tx.Commit()

	ctx.JSON(http.StatusOK, response.Response{
		Code:    20000,
		Message: "Success",
	})

}

func DoPackageByUpdate(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: "打包ID获取失败",
		})
		return
	}

	var former asset.DoPackageByUpdateForm
	if err := ctx.ShouldBind(&former); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: err.Error(),
		})
		return
	}

	var pack model.DorPackage
	data.Database.Preload("Details").First(&pack, id)
	if pack.Id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "未找到该打包数据",
		})
		return
	}

	var deviceIds = make([]uint, len(former.Devices))

	for idx, item := range former.Devices {
		deviceIds[idx] = item.Device
	}

	var devices []model.DorDevice
	data.Database.Where("id in (?)", deviceIds).Find(&devices)

	for _, item := range former.Devices {
		mark := true
		for _, value := range devices {
			if item.Device == value.Id {
				mark = false
			}
		}
		if mark {
			ctx.JSON(http.StatusOK, response.Response{
				Code:    40400,
				Message: "部分设备未找到",
			})
			return
		}
	}

	if len(devices) != len(former.Devices) {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "部分设备选择重复",
		})
		return
	}

	var creates, updates []model.DorPackageDetail
	for _, item := range former.Devices {
		mark := true
		for _, value := range pack.Details {
			if value.DeviceId == item.Device {
				mark = false
				if item.Number != value.Number {
					value.Number = item.Number
					updates = append(updates, value)
				}
			}
		}
		if mark {
			creates = append(creates, model.DorPackageDetail{
				PackageId: pack.Id,
				DeviceId:  item.Device,
				Number:    item.Number,
			})
		}
	}

	var deletes []uint
	for _, item := range pack.Details {
		mark := true
		for _, value := range former.Devices {
			if item.DeviceId == value.Device {
				mark = false
			}
		}
		if mark {
			deletes = append(deletes, item.Id)
		}
	}

	tx := data.Database.Begin()

	if former.Name != pack.Name {
		pack.Name = former.Name

		if t := tx.Save(&pack); t.RowsAffected <= 0 {
			tx.Rollback()
			ctx.JSON(http.StatusOK, response.Response{
				Code:    60000,
				Message: "修改失败",
			})
			return
		}
	}

	if len(creates) > 0 {
		if t := tx.Save(&creates); t.RowsAffected <= 0 {
			tx.Rollback()
			ctx.JSON(http.StatusOK, response.Response{
				Code:    60000,
				Message: "修改失败",
			})
			return
		}
	}
	if len(updates) > 0 {
		for _, item := range updates {
			if t := tx.Save(&item); t.RowsAffected <= 0 {
				tx.Rollback()
				ctx.JSON(http.StatusOK, response.Response{
					Code:    60000,
					Message: "修改失败",
				})
				return
			}
		}
	}
	if len(deletes) > 0 {
		if t := tx.Delete(&model.DorPackageDetail{}, deletes); t.RowsAffected <= 0 {
			tx.Rollback()
			ctx.JSON(http.StatusOK, response.Response{
				Code:    60000,
				Message: "修改失败",
			})
			return
		}
	}

	tx.Commit()

	ctx.JSON(http.StatusOK, response.Response{
		Code:    20000,
		Message: "Success",
	})

}

func DoPackageByDelete(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: "打包ID获取失败",
		})
		return
	}

	var pack model.DorPackage
	data.Database.First(&pack, id)
	if pack.Id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "未找到该打包数据",
		})
		return
	}

	tx := data.Database.Begin()

	if t := tx.Delete(&pack); t.RowsAffected <= 0 {
		tx.Rollback()
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: "删除失败",
		})
		return
	}

	if t := tx.Where("package_id=?", pack.Id).Delete(&model.DorPackageDetail{}); t.RowsAffected <= 0 {
		tx.Rollback()
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: "删除失败",
		})
		return
	}

	tx.Commit()

	ctx.JSON(http.StatusOK, response.Response{
		Code:    20000,
		Message: "Success",
	})

}

func ToPackageByPaginate(ctx *gin.Context) {

	var query asset.ToPackageByPaginateForm
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: err.Error(),
		})
		return
	}

	responses := response.Paginate{
		Code:    20000,
		Message: "Success",
	}

	responses.Data.Page = query.GetPage()
	responses.Data.Size = query.GetSize()
	responses.Data.Data = []any{}

	tx := data.Database

	if query.Keyword != "" {
		tx = tx.Where("name like ?", "%"+query.Keyword+"%")
	}

	tc := tx
	tc.Table(model.TableDorPackage).Count(&responses.Data.Total)

	if responses.Data.Total > 0 {

		var packages []model.DorPackage

		tx.Preload("Details").Preload("Details.Device").Offset(query.GetOffset()).Limit(query.GetLimit()).Find(&packages)

		for _, item := range packages {
			items := assetResponse.ToPackageByPaginateResponse{
				Id:        item.Id,
				Name:      item.Name,
				CreatedAt: item.CreatedAt.ToDateTimeString(),
			}
			for _, value := range item.Details {
				items.Devices = append(items.Devices, assetResponse.ToPackageByPaginateOfDevicesResponse{
					Id:       value.DeviceId,
					Category: value.Device.CategoryId,
					Name:     value.Device.Name,
					Number:   value.Number,
				})
			}
			responses.Data.Data = append(responses.Data.Data, items)
		}
	}

	ctx.JSON(http.StatusOK, responses)

}

func ToPackageByOnline(ctx *gin.Context) {

	responses := response.Responses{
		Code:    20000,
		Message: "Success",
		Data:    []any{},
	}

	var packages []model.DorPackage

	data.Database.Order("`id` desc").Find(&packages)
	for _, item := range packages {
		responses.Data = append(responses.Data, assetResponse.ToPackageByOnlineResponse{
			Id:   item.Id,
			Name: item.Name,
		})
	}

	ctx.JSON(http.StatusOK, responses)
}
