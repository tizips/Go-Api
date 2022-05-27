package asset

import (
	"github.com/gin-gonic/gin"
	"saas/app/model"
	"saas/app/request/admin/dormitory/asset"
	assetResponse "saas/app/response/admin/dormitory/asset"
	"saas/kernel/data"
	"saas/kernel/response"
	"strconv"
)

func DoPackageByCreate(ctx *gin.Context) {

	var request asset.DoPackageByCreate
	if err := ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	var deviceIds = make([]uint, len(request.Devices))

	for idx, item := range request.Devices {
		deviceIds[idx] = item.Device
	}

	var devices []model.DorDevice
	data.Database.Where("`id` in (?)", deviceIds).Find(&devices)

	for _, item := range request.Devices {
		mark := true
		for _, value := range devices {
			if item.Device == value.Id {
				mark = false
			}
		}
		if mark {
			response.NotFound(ctx, "部分设备未找到")
			return
		}
	}

	if len(devices) != len(request.Devices) {
		response.Fail(ctx, "部分设备选择重复")
		return
	}

	tx := data.Database.Begin()

	pack := model.DorPackage{
		Name: request.Name,
	}

	if t := tx.Create(&pack); t.RowsAffected <= 0 {
		tx.Rollback()
		response.Fail(ctx, "添加失败")
		return
	}

	var bindings []model.DorPackageDetail
	for _, item := range request.Devices {
		bindings = append(bindings, model.DorPackageDetail{
			PackageId: pack.Id,
			DeviceId:  item.Device,
			Number:    item.Number,
		})
	}

	if t := tx.Create(&bindings); t.RowsAffected <= 0 {
		tx.Rollback()
		response.Fail(ctx, "添加失败")
		return
	}

	tx.Commit()

	response.Success(ctx)
}

func DoPackageByUpdate(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		response.FailByRequestWithMessage(ctx, "ID获取失败")
		return
	}

	var request asset.DoPackageByUpdate
	if err = ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	var pack model.DorPackage
	data.Database.Preload("Details").Find(&pack, id)
	if pack.Id <= 0 {
		response.NotFound(ctx, "未找到该打包数据")
		return
	}

	var deviceIds = make([]uint, len(request.Devices))

	for idx, item := range request.Devices {
		deviceIds[idx] = item.Device
	}

	var devices []model.DorDevice
	data.Database.Where("id in (?)", deviceIds).Find(&devices)

	for _, item := range request.Devices {
		mark := true
		for _, value := range devices {
			if item.Device == value.Id {
				mark = false
			}
		}
		if mark {
			response.NotFound(ctx, "部分设备未找到")
			return
		}
	}

	if len(devices) != len(request.Devices) {
		response.Fail(ctx, "部分设备选择重复")
		return
	}

	var creates, updates []model.DorPackageDetail
	for _, item := range request.Devices {
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
		for _, value := range request.Devices {
			if item.DeviceId == value.Device {
				mark = false
			}
		}
		if mark {
			deletes = append(deletes, item.Id)
		}
	}

	tx := data.Database.Begin()

	if request.Name != pack.Name {
		pack.Name = request.Name

		if t := tx.Save(&pack); t.RowsAffected <= 0 {
			tx.Rollback()
			response.Fail(ctx, "修改失败")
			return
		}
	}

	if len(creates) > 0 {
		if t := tx.Save(&creates); t.RowsAffected <= 0 {
			tx.Rollback()
			response.Fail(ctx, "修改失败")
			return
		}
	}
	if len(updates) > 0 {
		for _, item := range updates {
			if t := tx.Save(&item); t.RowsAffected <= 0 {
				tx.Rollback()
				response.Fail(ctx, "修改失败")
				return
			}
		}
	}
	if len(deletes) > 0 {
		if t := tx.Delete(&model.DorPackageDetail{}, deletes); t.RowsAffected <= 0 {
			tx.Rollback()
			response.Fail(ctx, "修改失败")
			return
		}
	}

	tx.Commit()

	response.Success(ctx)
}

func DoPackageByDelete(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		response.FailByRequestWithMessage(ctx, "ID获取失败")
		return
	}

	var pack model.DorPackage
	data.Database.Find(&pack, id)
	if pack.Id <= 0 {
		response.NotFound(ctx, "未找到该打包数据")
		return
	}

	tx := data.Database.Begin()

	if t := tx.Delete(&pack); t.RowsAffected <= 0 {
		tx.Rollback()
		response.Fail(ctx, "删除失败")
		return
	}

	if t := tx.Where("`package_id`=?", pack.Id).Delete(&model.DorPackageDetail{}); t.RowsAffected <= 0 {
		tx.Rollback()
		response.Fail(ctx, "删除失败")
		return
	}

	tx.Commit()

	response.Success(ctx)
}

func ToPackageByPaginate(ctx *gin.Context) {

	var request asset.ToPackageByPaginate
	if err := ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	responses := response.Paginate{
		Total: 0,
		Page:  request.GetPage(),
		Size:  request.GetSize(),
		Data:  make([]any, 0),
	}

	tx := data.Database

	if request.Keyword != "" {
		tx = tx.Where("`name` like ?", "%"+request.Keyword+"%")
	}

	tc := tx
	tc.Model(&model.DorPackage{}).Count(&responses.Total)

	if responses.Total > 0 {

		var packages []model.DorPackage

		tx.Preload("Details.Device").Offset(request.GetOffset()).Limit(request.GetLimit()).Find(&packages)

		for _, item := range packages {
			items := assetResponse.ToPackageByPaginate{
				Id:        item.Id,
				Name:      item.Name,
				CreatedAt: item.CreatedAt.ToDateTimeString(),
			}
			for _, value := range item.Details {
				items.Devices = append(items.Devices, assetResponse.ToPackageByPaginateOfDevices{
					Id:       value.DeviceId,
					Category: value.Device.CategoryId,
					Name:     value.Device.Name,
					Number:   value.Number,
				})
			}
			responses.Data = append(responses.Data, items)
		}
	}

	response.SuccessByPaginate(ctx, responses)
}

func ToPackageByOnline(ctx *gin.Context) {

	responses := make([]any, 0)

	var packages []model.DorPackage

	data.Database.Order("`id` desc").Find(&packages)
	for _, item := range packages {
		responses = append(responses, assetResponse.ToPackageByOnline{
			Id:   item.Id,
			Name: item.Name,
		})
	}

	response.SuccessByList(ctx, responses)
}
