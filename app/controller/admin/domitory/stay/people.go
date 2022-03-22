package stay

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
	"net/http"
	"saas/app/constant"
	"saas/app/form/admin/dormitory/stay"
	"saas/app/model"
	stayResponse "saas/app/response/admin/dormitory/stay"
	"saas/kernel/data"
	"saas/kernel/response"
)

func ToPeopleByPaginate(ctx *gin.Context) {

	var query stay.ToPeopleByPaginateForm
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: err.Error(),
		})
		return
	}

	tx := data.Database.Where(fmt.Sprintf("%s.status=?", model.TableDorPeople), query.Status)

	if query.Floor > 0 {
		tx = tx.Where("floor_id=?", query.Floor)
	} else if query.Building > 0 {
		tx = tx.Where("building_id=?", query.Building)
	}

	if query.IsTemp > 0 {
		tx = tx.Where(fmt.Sprintf("%s.is_temp=?", model.TableDorPeople), query.IsTemp)
	}

	if query.Keyword != "" {

		condition := data.Database.Select("1")

		if query.Type == "mobile" {
			condition = condition.
				Table(model.TableMemMember).
				Where(fmt.Sprintf("%s.member_id=%s.id", model.TableDorPeople, model.TableMemMember)).
				Where(fmt.Sprintf("%s.mobile=?", model.TableMemMember), query.Keyword)
		} else if query.Type == "room" {
			condition = condition.
				Table(model.TableDorRoom).
				Where(fmt.Sprintf("%s.room_id=%s.id", model.TableDorPeople, model.TableDorRoom)).
				Where(fmt.Sprintf("%s.name=?", model.TableDorRoom), query.Keyword)
		} else {
			condition = condition.
				Table(model.TableMemMember).
				Where(fmt.Sprintf("%s.member_id=%s.id", model.TableDorPeople, model.TableMemMember)).
				Where(fmt.Sprintf("%s.name=?", model.TableMemMember), query.Keyword)
		}

		tx = tx.Where("exists (?)", condition)
	}

	tc := tx

	responses := response.Paginate{
		Code:    20000,
		Message: "Success",
	}

	responses.Data.Size = query.GetSize()
	responses.Data.Page = query.GetPage()
	responses.Data.Data = []any{}

	tc.Model(model.DorPeople{}).Count(&responses.Data.Total)

	if responses.Data.Total > 0 {
		var peoples []model.DorPeople

		tx.
			Preload("Member", func(db *gorm.DB) *gorm.DB { return db.Unscoped() }).
			Preload("Staff", func(db *gorm.DB) *gorm.DB { return db.Unscoped() }).
			Preload("Certification", func(db *gorm.DB) *gorm.DB { return db.Unscoped() }).
			Preload("Category", func(db *gorm.DB) *gorm.DB { return db.Unscoped() }).
			Preload("Building", func(db *gorm.DB) *gorm.DB { return db.Unscoped() }).
			Preload("Floor", func(db *gorm.DB) *gorm.DB { return db.Unscoped() }).
			Preload("Room", func(db *gorm.DB) *gorm.DB { return db.Unscoped() }).
			Preload("Bed", func(db *gorm.DB) *gorm.DB { return db.Unscoped() }).
			Order(fmt.Sprintf("%s.`id` desc", model.TableDorPeople)).
			Offset(query.GetOffset()).
			Limit(query.GetLimit()).
			Find(&peoples)

		for _, item := range peoples {

			results := stayResponse.ToPeopleByPaginateResponse{
				Id:        item.Id,
				Category:  item.Category.Name,
				Building:  item.Building.Name,
				Floor:     item.Floor.Name,
				Room:      item.Room.Name,
				Bed:       item.Bed.Name,
				Name:      item.Member.Name,
				Mobile:    item.Member.Mobile,
				IsTemp:    item.Category.IsTemp,
				Start:     item.Start.ToDateString(),
				Remark:    item.Remark,
				CreatedAt: item.CreatedAt.ToDateTimeString(),
			}
			if item.End != nil {
				results.End = item.End.ToDateTimeString()
			}
			if item.Staff != nil && item.Staff.Id > 0 {
				results.Staff = item.Staff.Status
				results.Titles = item.Staff.Title
			}
			if item.Certification != nil && item.Certification.Id > 0 {
				certification := stayResponse.ToPeopleByPaginateOfCertificationResponse{
					No: item.Certification.No,
				}
				results.Certification = certification
			}
			responses.Data.Data = append(responses.Data.Data, results)
		}
	}

	ctx.JSON(http.StatusOK, responses)

}

func DoPeopleByCreate(ctx *gin.Context) {

	var former stay.DoPeopleByCreateForm
	if err := ctx.ShouldBind(&former); err != nil {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40000,
			Message: err.Error(),
		})
		return
	}

	var bed model.DorBed
	data.Database.Preload("Building").Preload("Floor").Preload("Room").Preload("Type").Where("is_enable=?", constant.IsEnableYes).First(&bed, former.Bed)
	if bed.Id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "床位不存在",
		})
		return
	}

	var category model.DorStayCategory
	data.Database.Where("is_enable", constant.IsEnableYes).First(&category, former.Category)
	if category.Id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "类型不存在",
		})
		return
	}

	var count int64 = 0
	data.Database.Model(model.DorPeople{}).Joins(fmt.Sprintf("left join %s on %s.member_id=%s.id", model.TableMemMember, model.TableDorPeople, model.TableMemMember)).Where(fmt.Sprintf("%s.mobile=?", model.TableMemMember), former.Mobile).Where(fmt.Sprintf("%s.status=?", model.TableDorPeople), model.DorPeopleStatusLive).Count(&count)
	if count > 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "该手机号已办理了入住，无法重复办理",
		})
		return
	}

	var member model.MemMember
	data.Database.Where("mobile=?", former.Mobile).First(&member)
	if member.Id == "" {

		node, err := snowflake.NewNode(1)
		if err != nil {
			ctx.JSON(http.StatusOK, response.Response{
				Code:    40400,
				Message: "办理失败",
			})
			return
		}

		member = model.MemMember{
			Id:       node.Generate().String(),
			Mobile:   former.Mobile,
			Name:     former.Name,
			Nickname: former.Name,
			IsEnable: constant.IsEnableYes,
		}

		if t := data.Database.Create(&member); t.RowsAffected <= 0 {
			ctx.JSON(http.StatusOK, response.Response{
				Code:    60000,
				Message: "办理失败",
			})
			return
		}
	}

	tx := data.Database.Begin()

	people := model.DorPeople{
		CategoryId: category.Id,
		BuildingId: bed.BuildingId,
		FloorId:    bed.FloorId,
		RoomId:     bed.RoomId,
		BedId:      bed.Id,
		TypeId:     bed.TypeId,
		MemberId:   member.Id,
		Start:      carbon.Date{Carbon: carbon.Parse(former.Start)},
		Status:     model.DorPeopleStatusLive,
		IsTemp:     category.IsTemp,
		Remark:     former.Remark,
	}

	if former.End != "" {
		people.End = &carbon.Date{Carbon: carbon.Parse(former.End)}
	}

	if t := tx.Create(&people); t.RowsAffected <= 0 {
		tx.Rollback()
		ctx.JSON(http.StatusOK, response.Response{
			Code:    60000,
			Message: "办理失败",
		})
		return
	}

	var masterId uint = 0
	masterName := member.Name

	var master model.DorPeople
	tx.Preload("Master.Member").Where("bed_id=?", people.BedId).Where("master_id<>?", 0).Where("status=?", model.DorPeopleStatusLive).First(&master)
	if master.MasterId > 0 {
		masterId = master.MasterId
		masterName = master.Member.Name
	} else {
		masterId = people.Id
	}

	people.MasterId = masterId

	if t := tx.Table(model.TableDorPeople).Where("id=?", people.Id).UpdateColumn("master_id", masterId); t.RowsAffected <= 0 {
		tx.Rollback()
		ctx.JSON(http.StatusOK, response.Response{
			Code:    60000,
			Message: "办理失败",
		})
		return
	}

	var details = make(map[string]any, 11)

	details["category"] = category.Name
	details["building"] = bed.Building.Name
	details["floor"] = bed.Floor.Name
	details["room"] = bed.Room.Name
	details["bed"] = bed.Name
	details["type"] = bed.Type.Name
	details["name"] = member.Name
	details["mobile"] = member.Mobile
	details["master"] = masterName
	details["is_temp"] = category.IsTemp
	details["start"] = people.Start
	details["end"] = people.End

	str, _ := json.Marshal(details)

	log := model.DorPeopleLog{
		PeopleId: people.Id,
		MemberId: member.Id,
		Status:   model.DorPeopleLogStatusLive,
		Detail:   string(str),
		Remark:   people.Remark,
	}

	if t := tx.Create(&log); t.RowsAffected <= 0 {
		tx.Rollback()
		ctx.JSON(http.StatusOK, response.Response{
			Code:    60000,
			Message: "办理失败",
		})
		return
	}

	tx.Commit()

	ctx.JSON(http.StatusOK, response.Response{
		Code:    20000,
		Message: "Success",
	})
}
