package asset

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon/v2"
	"github.com/gookit/goutil/dump"
	"gorm.io/gorm"
	"saas/app/form/admin/dormitory/asset"
	"saas/app/helper/collection"
	"saas/app/model"
	assetResponse "saas/app/response/admin/dormitory/asset"
	"saas/kernel/data"
	"saas/kernel/response"
)

func DoGrantByCreate(ctx *gin.Context) {

	var former asset.DoGrantByCreateFormer
	if err := ctx.ShouldBind(&former); err != nil {
		response.ToResponseByFailRequest(ctx, err)
		return
	}

	var pack model.DorPackage
	var device model.DorDevice

	if former.Package > 0 {
		data.Database.Preload("Details").Find(&pack, former.Package)
		if pack.Id <= 0 {
			response.ToResponseByNotFound(ctx, "打包不存在")
			return
		}
	} else if former.Device > 0 {
		data.Database.Find(&device, former.Device)
		if device.Id <= 0 {
			response.ToResponseByNotFound(ctx, "设备不存在")
			return
		}
	}

	var typeBed model.DorTypeBed
	var buildingIds, floorIds, roomIds, bedIds []uint

	if former.Type > 0 {
		data.Database.Find(&typeBed, former.Type)
		if typeBed.Id <= 0 {
			response.ToResponseByNotFound(ctx, "房型位置不存在")
			return
		}
	} else {
		for _, item := range former.Positions {
			if item.Object == "bed" {
				bedIds = append(bedIds, item.Id)
			} else if item.Object == "room" {
				roomIds = append(roomIds, item.Id)
			} else if item.Object == "floor" {
				floorIds = append(floorIds, item.Id)
			} else if item.Object == "building" {
				buildingIds = append(buildingIds, item.Id)
			}
		}
	}

	dump.P(former)
	dump.P(buildingIds)

	tx := data.Database.Begin()

	grant := model.DorGrant{
		Object:    "device",
		PackageId: pack.Id,
		Remark:    former.Remark,
	}
	if pack.Id > 0 {
		grant.Object = model.DorGrantObjectPackage
	}
	if t := tx.Create(&grant); t.RowsAffected <= 0 {
		tx.Rollback()
		response.ToResponseByFail(ctx, "发放失败")
		return
	}

	var devices []model.DorGrantDevice
	if pack.Id > 0 {
		for _, item := range pack.Details {
			devices = append(devices, model.DorGrantDevice{
				GrantId:  grant.Id,
				DeviceId: item.DeviceId,
				Number:   item.Number,
			})
		}
	} else if device.Id > 0 {
		devices = append(devices, model.DorGrantDevice{
			GrantId:  grant.Id,
			DeviceId: device.Id,
			Number:   former.Number,
		})
	}
	if t := tx.Create(&devices); t.RowsAffected <= 0 {
		tx.Rollback()
		response.ToResponseByFail(ctx, "发放失败")
		return
	}

	var positions []model.DorGrantPosition
	var details []model.DorGrantDetail

	var isNoPublicBuildingIds []uint
	var isNoPublicFloorIds []uint
	var isNoPublicRoomIds []uint

	if former.Type > 0 {
		positions = append(positions, model.DorGrantPosition{
			GrantId:   grant.Id,
			Object:    model.DorGrantPositionType,
			TypeId:    typeBed.TypeId,
			TypeBedId: typeBed.Id,
		})
	} else {
		if len(buildingIds) > 0 {
			var buildings []model.DorBuilding
			handleBuildingIds := collection.Unique(buildingIds)
			if len(handleBuildingIds) != len(buildingIds) {
				response.ToResponseByFail(ctx, "楼栋选择重复")
				return
			}
			data.Database.Find(&buildings, buildingIds)
			if len(buildings) <= 0 {
				response.ToResponseByNotFound(ctx, "楼栋未找到")
				return
			} else if len(buildings) != len(buildingIds) {
				response.ToResponseByNotFound(ctx, "部分楼栋未找到")
				return
			}
			for _, item := range buildings {
				positions = append(positions, model.DorGrantPosition{
					GrantId:    grant.Id,
					Object:     model.DorGrantPositionBuilding,
					BuildingId: item.Id,
				})
				if item.IsPublic != model.DorBuildingIsPublicYes {
					isNoPublicBuildingIds = append(isNoPublicBuildingIds, item.Id)
				} else {
					for _, value := range devices {
						details = append(details, model.DorGrantDetail{
							GrantId:    grant.Id,
							PackageId:  pack.Id,
							TypeId:     typeBed.TypeId,
							BuildingId: item.Id,
							DeviceId:   value.DeviceId,
							Number:     value.Number,
							IsPublic:   model.DorBuildingIsPublicYes,
						})
					}
				}
			}
		}
		if len(floorIds) > 0 {
			var floors []model.DorFloor
			handleFloorIds := collection.Unique(floorIds)
			if len(handleFloorIds) != len(floorIds) {
				response.ToResponseByFail(ctx, "楼层选择重复")
				return
			}
			data.Database.Find(&floors, floorIds)
			if len(floors) <= 0 {
				response.ToResponseByNotFound(ctx, "楼层未找到")
				return
			} else if len(floors) != len(floorIds) {
				response.ToResponseByNotFound(ctx, "部分楼层未找到")
				return
			}
			for _, item := range floors {
				positions = append(positions, model.DorGrantPosition{
					GrantId:    grant.Id,
					Object:     model.DorGrantPositionFloor,
					BuildingId: item.BuildingId,
					FloorId:    item.Id,
				})
				if item.IsPublic != model.DorFloorIsPublicYes {
					isNoPublicFloorIds = append(isNoPublicFloorIds, item.Id)
				} else {
					for _, value := range devices {
						details = append(details, model.DorGrantDetail{
							GrantId:    grant.Id,
							PackageId:  pack.Id,
							TypeId:     typeBed.TypeId,
							BuildingId: item.BuildingId,
							FloorId:    item.Id,
							DeviceId:   value.DeviceId,
							Number:     value.Number,
							IsPublic:   model.DorBuildingIsPublicYes,
						})
					}
				}
			}
		}
		if len(roomIds) > 0 {
			var rooms []model.DorRoom
			handleRoomIds := collection.Unique(roomIds)
			if len(handleRoomIds) != len(roomIds) {
				response.ToResponseByFail(ctx, "房间选择重复")
				return
			}
			data.Database.Find(&rooms, roomIds)
			if len(rooms) <= 0 {
				response.ToResponseByNotFound(ctx, "房间未找到")
				return
			} else if len(rooms) != len(roomIds) {
				response.ToResponseByNotFound(ctx, "部分房间未找到")
				return
			}
			for _, item := range rooms {
				positions = append(positions, model.DorGrantPosition{
					GrantId:    grant.Id,
					Object:     model.DorGrantPositionRoom,
					BuildingId: item.BuildingId,
					FloorId:    item.FloorId,
					RoomId:     item.Id,
				})
				if item.IsPublic != model.DorRoomIsPublicYes {
					isNoPublicRoomIds = append(isNoPublicRoomIds, item.Id)
				} else {
					for _, value := range devices {
						details = append(details, model.DorGrantDetail{
							GrantId:    grant.Id,
							PackageId:  pack.Id,
							TypeId:     typeBed.TypeId,
							BuildingId: item.BuildingId,
							FloorId:    item.FloorId,
							RoomId:     item.Id,
							DeviceId:   value.DeviceId,
							Number:     value.Number,
							IsPublic:   model.DorBuildingIsPublicYes,
						})
					}
				}
			}
		}
		if len(bedIds) > 0 {
			var beds []model.DorBed
			handleBedIds := collection.Unique(bedIds)
			if len(handleBedIds) != len(bedIds) {
				response.ToResponseByFail(ctx, "床位选择重复")
				return
			}
			data.Database.Find(&beds, bedIds)
			if len(beds) <= 0 {
				response.ToResponseByNotFound(ctx, "床位未找到")
				return
			} else if len(beds) != len(bedIds) {
				response.ToResponseByNotFound(ctx, "部分床位未找到")
				return
			}
			for _, item := range beds {
				positions = append(positions, model.DorGrantPosition{
					GrantId:    grant.Id,
					Object:     model.DorGrantPositionBed,
					BuildingId: item.BuildingId,
					FloorId:    item.FloorId,
					RoomId:     item.RoomId,
					BedId:      item.Id,
				})
			}
		}
	}
	if t := tx.Create(&positions); t.RowsAffected <= 0 {
		tx.Rollback()
		response.ToResponseByFail(ctx, "发放失败")
		return
	}

	if former.Type > 0 || len(isNoPublicBuildingIds) > 0 || len(isNoPublicFloorIds) > 0 || len(isNoPublicRoomIds) > 0 || len(bedIds) > 0 {

		var results []map[string]any

		tb := data.Database.
			Select(fmt.Sprintf("%s.*, %s.id as people_id, %s.id as member_id", model.TableDorBed, model.TableDorPeople, model.TableDorPeople)).
			Table(model.TableDorBed).
			Joins(fmt.Sprintf("left join %s on %s.id=%s.bed_id and %s.status", model.TableDorPeople, model.TableDorBed, model.TableDorPeople, model.TableDorPeople))

		if former.Type > 0 {
			tb = tb.Where(fmt.Sprintf("%s.bed_id=?", model.TableDorBed), typeBed.Id)
		} else {
			condition := data.Database

			if len(isNoPublicBuildingIds) > 0 {
				condition = condition.
					Or(data.Database.
						Where(fmt.Sprintf("%s.is_public=?", model.TableDorBed), model.DorBedIsPublicNo).
						Where(fmt.Sprintf("%s.building_id in ?", model.TableDorBed), isNoPublicBuildingIds),
					)
			}
			if len(isNoPublicFloorIds) > 0 {
				condition = condition.
					Or(data.Database.
						Where(fmt.Sprintf("%s.is_public=?", model.TableDorBed), model.DorBedIsPublicNo).
						Where(fmt.Sprintf("%s.floor_id_id in ?", model.TableDorBed), isNoPublicFloorIds),
					)
			}
			if len(isNoPublicRoomIds) > 0 {
				condition = condition.
					Or(data.Database.
						Where(fmt.Sprintf("%s.is_public=?", model.TableDorBed), model.DorBedIsPublicNo).
						Where(fmt.Sprintf("%s.room_id in ?", model.TableDorBed), isNoPublicRoomIds),
					)
			}
			if len(bedIds) > 0 {
				condition = condition.Or(fmt.Sprintf("%s.id in ?", model.TableDorBed), bedIds)
			}

			tb = tb.Where(condition)
		}

		tb.Find(&results)

		if len(results) <= 0 {
			tx.Rollback()
			response.ToResponseByFail(ctx, "该位置尚未配备具体床位")
			return
		}

		for _, item := range results {
			for _, value := range devices {
				items := model.DorGrantDetail{
					GrantId:    grant.Id,
					PackageId:  grant.PackageId,
					PositionId: 0,
					TypeId:     uint(item["type_id"].(uint32)),
					BuildingId: uint(item["building_id"].(uint32)),
					FloorId:    uint(item["floor_id"].(uint32)),
					RoomId:     uint(item["room_id"].(uint32)),
					BedId:      uint(item["id"].(uint32)),
					DeviceId:   value.DeviceId,
					Number:     value.Number,
					IsPublic:   item["is_public"].(uint8),
				}
				if item["people_id"] != nil {
					items.PeopleId = uint(item["people_id"].(uint32))
				}
				if item["member_id"] != nil {
					items.MemberId = item["member_id"].(string)
				}
				details = append(details, items)
			}
		}
	}

	if len(details) > 0 {
		for index, item := range details {
			for _, val := range positions {
				if val.Object == model.DorGrantPositionType && item.TypeId == val.TypeId {
					details[index].PositionId = val.Id
					break
				} else if val.Object == model.DorGrantPositionBuilding && item.BuildingId == val.BuildingId {
					details[index].PositionId = val.Id
					break
				} else if val.Object == model.DorGrantPositionFloor && item.FloorId == val.FloorId {
					details[index].PositionId = val.Id
					break
				} else if val.Object == model.DorGrantPositionRoom && item.RoomId == val.RoomId {
					details[index].PositionId = val.Id
					break
				} else if val.Object == model.DorGrantPositionBed && item.BedId == val.BedId {
					details[index].PositionId = val.Id
					break
				}
			}
		}
		if t := tx.CreateInBatches(&details, 20); t.RowsAffected <= 0 {
			tx.Rollback()
			response.ToResponseByFail(ctx, "发放失败")
			return
		}
	}

	tx.Commit()

	response.ToResponseBySuccess(ctx)
}

func DoGrantByRevoke(ctx *gin.Context) {

	var former asset.DoGrantByRevokeFormer
	if err := ctx.ShouldBind(&former); err != nil {
		response.ToResponseByFailRequest(ctx, err)
		return
	}

	var grant model.DorGrant
	data.Database.First(&grant, former.Id)
	if grant.Id <= 0 {
		response.ToResponseByNotFound(ctx, "发放记录不存在")
		return
	}

	hours := grant.CreatedAt.DiffInSecondsWithAbs(carbon.Now())
	if hours > 86400 {
		response.ToResponseByFail(ctx, "本次发放已无法撤销")
		return
	}

	tx := data.Database

	//if t := tx.Save(&grant); t.RowsAffected <= 0 {
	//	ctx.JSON(http.StatusOK, response.Response{
	//		Code:    40400,
	//		Message: "回撤失败",
	//	})
	//	return
	//}

	if t := tx.Delete(&grant); t.RowsAffected <= 0 {
		response.ToResponseByFail(ctx, "回撤失败")
		return
	}

	if t := tx.Where("grant_id=?", grant.Id).Delete(&model.DorGrantPosition{}); t.RowsAffected <= 0 {
		response.ToResponseByFail(ctx, "回撤失败")
		return
	}

	if t := tx.Where("grant_id=?", grant.Id).Delete(&model.DorGrantDevice{}); t.RowsAffected <= 0 {
		response.ToResponseByFail(ctx, "回撤失败")
		return
	}

	if t := tx.Where("grant_id=?", grant.Id).Delete(&model.DorGrantDetail{}); t.RowsAffected <= 0 {
		response.ToResponseByFail(ctx, "回撤失败")
		return
	}

	tx.Commit()

	response.ToResponseBySuccess(ctx)
}

func ToGrantByPaginate(ctx *gin.Context) {

	var query asset.ToGrantByPaginateFormer
	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.ToResponseByFailRequest(ctx, err)
		return
	}

	responses := response.Paginate{
		Page: query.GetPage(),
		Size: query.GetSize(),
		Data: make([]any, 0),
	}

	tx := data.Database

	tc := tx

	tc.Table(model.TableDorGrant).Count(&responses.Total)

	if responses.Total > 0 {
		var grants []model.DorGrant
		tx.
			Preload("Package", func(db *gorm.DB) *gorm.DB { return db.Unscoped() }).
			Preload("BindDevices", func(db *gorm.DB) *gorm.DB { return db.Unscoped() }).
			Preload("BindDevices.Device", func(db *gorm.DB) *gorm.DB { return db.Unscoped() }).
			Order("id desc").
			Offset(query.GetOffset()).
			Limit(query.GetLimit()).
			Find(&grants)

		for _, item := range grants {
			items := assetResponse.ToGrantByPaginateResponse{
				Id:        item.Id,
				Remark:    item.Remark,
				CreatedAt: item.CreatedAt.ToDateTimeString(),
			}
			if item.Object == model.DorGrantObjectPackage {
				items.Package = item.Package.Name
			}
			for _, value := range item.BindDevices {
				items.Devices = append(items.Devices, assetResponse.ToGrantByPaginateOfDeviceResponse{
					Name:   value.Device.Name,
					Number: value.Number,
				})
			}

			responses.Data = append(responses.Data, items)
		}
	}

	response.ToResponseBySuccessPaginate(ctx, responses)
}
