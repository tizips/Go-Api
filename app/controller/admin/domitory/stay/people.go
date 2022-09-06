package stay

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"replace.github.com/bsm/redislock"
	"saas/app/constant"
	"saas/app/model"
	"saas/app/request/admin/dormitory/stay"
	res "saas/app/response/admin/dormitory/stay"
	"saas/kernel/app"
	"saas/kernel/response"
	"time"
)

func ToPeopleByPaginate(ctx *gin.Context) {

	var request stay.ToPeopleByPaginate

	if err := ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	tx := app.Database.Where(fmt.Sprintf("%s.`status`=?", model.TableDorPeople), request.Status)

	if request.Floor > 0 {
		tx = tx.Where("`floor_id`=?", request.Floor)
	} else if request.Building > 0 {
		tx = tx.Where("`building_id`=?", request.Building)
	}

	if request.IsTemp > 0 {
		tx = tx.Where(fmt.Sprintf("%s.`is_temp`=?", model.TableDorPeople), request.IsTemp)
	}

	if request.Keyword != "" {

		condition := app.Database.Select("1")

		if request.Type == "mobile" {
			condition = condition.
				Table(model.TableMemMember).
				Where(fmt.Sprintf("%s.`member_id`=%s.`id`", model.TableDorPeople, model.TableMemMember)).
				Where(fmt.Sprintf("%s.`mobile`=?", model.TableMemMember), request.Keyword)
		} else if request.Type == "room" {
			condition = condition.
				Table(model.TableDorRoom).
				Where(fmt.Sprintf("%s.`room_id`=%s.`id`", model.TableDorPeople, model.TableDorRoom)).
				Where(fmt.Sprintf("%s.`name`=?", model.TableDorRoom), request.Keyword)
		} else {
			condition = condition.
				Table(model.TableMemMember).
				Where(fmt.Sprintf("%s.`member_id`=%s.`id`", model.TableDorPeople, model.TableMemMember)).
				Where(fmt.Sprintf("%s.`name`=?", model.TableMemMember), request.Keyword)
		}

		tx = tx.Where("exists (?)", condition)
	}

	tc := tx

	responses := response.Paginate[res.ToPeopleByPaginate]{
		Total: 0,
		Page:  request.GetPage(),
		Size:  request.GetSize(),
		Data:  make([]res.ToPeopleByPaginate, 0),
	}

	tc.Model(&model.DorPeople{}).Count(&responses.Total)

	if responses.Total > 0 {

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
			Offset(request.GetOffset()).
			Limit(request.GetLimit()).
			Find(&peoples)

		for _, item := range peoples {

			results := res.ToPeopleByPaginate{
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

				certification := res.ToPeopleByPaginateOfCertification{
					No: item.Certification.No,
				}

				results.Certification = certification
			}

			responses.Data = append(responses.Data, results)
		}
	}

	response.SuccessByPaginate(ctx, responses)
}

func DoPeopleByCreate(ctx *gin.Context) {

	var request stay.DoPeopleByCreate

	if err := ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	var bed model.DorBed

	if app.Database.Preload(clause.Associations).Find(&bed, "`id`=? and `is_enable`=?", request.Bed, constant.IsEnableYes); bed.Id <= 0 {
		response.NotFound(ctx, "床位不存在")
		return
	}

	var category model.DorStayCategory

	if app.Database.Find(&category, "`id`=? and `is_enable`=?", request.Category, constant.IsEnableYes); category.Id <= 0 {
		response.NotFound(ctx, "类型不存在")
		return
	}

	var count int64 = 0

	app.Database.Model(model.DorPeople{}).Joins(fmt.Sprintf("left join `%s` on `%s`.`member_id`=`%s`.`id`", model.TableMemMember, model.TableDorPeople, model.TableMemMember)).Where(fmt.Sprintf("`%s`.`mobile`=? and `%s`.`status`=?", model.TableMemMember, model.TableDorPeople), request.Mobile, model.DorPeopleStatusLive).Count(&count)

	if count > 0 {
		response.Fail(ctx, "该手机号已办理了入住，无法重复办理")
		return
	}

	var member model.MemMember

	if app.Database.Find(&member, "`mobile`=?", request.Mobile); member.Id == "" {

		obtain, err := redislock.New(app.Redis).Obtain(ctx, "lock:member:"+request.Mobile, time.Second*30, &redislock.Options{RetryStrategy: redislock.LinearBackoff(time.Microsecond * 3)})

		if err != nil {
			response.Fail(ctx, "办理失败")
			return
		}

		defer obtain.Release(ctx)

		member = model.MemMember{
			Id:       app.Snowflake.Generate().String(),
			Mobile:   request.Mobile,
			Name:     request.Name,
			Nickname: request.Name,
			IsEnable: constant.IsEnableYes,
		}

		if t := app.Database.Create(&member); t.RowsAffected <= 0 {
			response.Fail(ctx, "办理失败")
			return
		}
	}

	tx := app.Database.Begin()

	people := model.DorPeople{
		CategoryId: category.Id,
		BuildingId: bed.BuildingId,
		FloorId:    bed.FloorId,
		RoomId:     bed.RoomId,
		BedId:      bed.Id,
		TypeId:     bed.TypeId,
		MemberId:   member.Id,
		Start:      carbon.Date{Carbon: carbon.Parse(request.Start)},
		Status:     model.DorPeopleStatusLive,
		IsTemp:     category.IsTemp,
		Remark:     request.Remark,
	}

	if request.End != "" {
		people.End = &carbon.Date{Carbon: carbon.Parse(request.End)}
	}

	if t := tx.Create(&people); t.RowsAffected <= 0 {
		tx.Rollback()
		response.Fail(ctx, "办理失败")
		return
	}

	var masterId int = 0

	masterName := member.Name

	var master model.DorPeople

	if tx.Preload("Master.Member").Find(&master, "`bed_id`=? and `master_id`<>? and `status`=?", people.BedId, 0, model.DorPeopleStatusLive); master.MasterId > 0 {
		masterId = master.MasterId
		masterName = master.Member.Name
	} else {
		masterId = people.Id
	}

	people.MasterId = masterId

	if t := tx.Model(&model.DorPeople{}).Where("`id`=?", people.Id).UpdateColumn("`master_id`=?", masterId); t.RowsAffected <= 0 {
		tx.Rollback()
		response.Fail(ctx, "办理失败")
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
		response.Fail(ctx, "办理失败")
		return
	}

	tx.Commit()

	response.Success(ctx)
}

func DoPeopleByLeave(ctx *gin.Context) {

	var request stay.DoPeopleByLeave

	if err := ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

}
