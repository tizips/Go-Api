package asset

import (
	"github.com/gin-gonic/gin"
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
		response.ToResponseByFailRequest(ctx, err)
		return
	}

	var deviceIds = make([]uint, len(former.Devices))

	for idx, item := range former.Devices {
		deviceIds[idx] = item.Device
	}

	var devices []model.DorDevice
	data.Database.Where("`id` in (?)", deviceIds).Find(&devices)

	for _, item := range former.Devices {
		mark := true
		for _, value := range devices {
			if item.Device == value.Id {
				mark = false
			}
		}
		if mark {
			response.ToResponseByNotFound(ctx, "部分设备未找到")
			return
		}
	}

	if len(devices) != len(former.Devices) {
		response.ToResponseByFail(ctx, "部分设备选择重复")
		return
	}

	tx := data.Database.Begin()

	pack := model.DorPackage{
		Name: former.Name,
	}

	if t := tx.Create(&pack); t.RowsAffected <= 0 {
		tx.Rollback()
		response.ToResponseByFail(ctx, "添加失败")
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
		response.ToResponseByFail(ctx, "添加失败")
		return
	}

	tx.Commit()

	response.ToResponseBySuccess(ctx)
}

func DoPackageByUpdate(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		response.ToResponseByFailRequestMessage(ctx, "ID获取失败")
		return
	}

	var former asset.DoPackageByUpdateForm
	if err = ctx.ShouldBind(&former); err != nil {
		response.ToResponseByFailRequest(ctx, err)
		return
	}

	var pack model.DorPackage
	data.Database.Preload("Details").First(&pack, id)
	if pack.Id <= 0 {
		response.ToResponseByNotFound(ctx, "未找到该打包数据")
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
			response.ToResponseByNotFound(ctx, "部分设备未找到")
			return
		}
	}

	if len(devices) != len(former.Devices) {
		response.ToResponseByFail(ctx, "部分设备选择重复")
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
			response.ToResponseByFail(ctx, "修改失败")
			return
		}
	}

	if len(creates) > 0 {
		if t := tx.Save(&creates); t.RowsAffected <= 0 {
			tx.Rollback()
			response.ToResponseByFail(ctx, "修改失败")
			return
		}
	}
	if len(updates) > 0 {
		for _, item := range updates {
			if t := tx.Save(&item); t.RowsAffected <= 0 {
				tx.Rollback()
				response.ToResponseByFail(ctx, "修改失败")
				return
			}
		}
	}
	if len(deletes) > 0 {
		if t := tx.Delete(&model.DorPackageDetail{}, deletes); t.RowsAffected <= 0 {
			tx.Rollback()
			response.ToResponseByFail(ctx, "修改失败")
			return
		}
	}

	tx.Commit()

	response.ToResponseBySuccess(ctx)
}

func DoPackageByDelete(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id <= 0 {
		response.ToResponseByFailRequestMessage(ctx, "ID获取失败")
		return
	}

	var pack model.DorPackage
	data.Database.First(&pack, id)
	if pack.Id <= 0 {
		response.ToResponseByNotFound(ctx, "未找到该打包数据")
		return
	}

	tx := data.Database.Begin()

	if t := tx.Delete(&pack); t.RowsAffected <= 0 {
		tx.Rollback()
		response.ToResponseByFail(ctx, "删除失败")
		return
	}

	if t := tx.Where("`package_id`=?", pack.Id).Delete(&model.DorPackageDetail{}); t.RowsAffected <= 0 {
		tx.Rollback()
		response.ToResponseByFail(ctx, "删除失败")
		return
	}

	tx.Commit()

	response.ToResponseBySuccess(ctx)
}

func ToPackageByPaginate(ctx *gin.Context) {

	var query asset.ToPackageByPaginateForm
	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.ToResponseByFailRequest(ctx, err)
		return
	}

	responses := response.Paginate{
		Total: 0,
		Page:  query.GetPage(),
		Size:  query.GetSize(),
		Data:  make([]any, 0),
	}

	tx := data.Database

	if query.Keyword != "" {
		tx = tx.Where("`name` like ?", "%"+query.Keyword+"%")
	}

	tc := tx
	tc.Model(&model.DorPackage{}).Count(&responses.Total)

	if responses.Total > 0 {

		var packages []model.DorPackage

		tx.Preload("Details.Device").Offset(query.GetOffset()).Limit(query.GetLimit()).Find(&packages)

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
			responses.Data = append(responses.Data, items)
		}
	}

	response.ToResponseBySuccessPaginate(ctx, responses)
}

func ToPackageByOnline(ctx *gin.Context) {

	responses := make([]any, 0)

	var packages []model.DorPackage

	data.Database.Order("`id` desc").Find(&packages)
	for _, item := range packages {
		responses = append(responses, assetResponse.ToPackageByOnlineResponse{
			Id:   item.Id,
			Name: item.Name,
		})
	}

	response.ToResponseBySuccessList(ctx, responses)
}
