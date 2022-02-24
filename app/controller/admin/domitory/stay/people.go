package stay

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
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

	tx := data.Database.Where(fmt.Sprintf("%s.status", model.TableDorPeople), query.Status)

	if query.Floor > 0 {
		tx.Where("floor_id", query.Building)
	} else if query.Building > 0 {
		tx.Where("building_id", query.Building)
	}

	if ctx.DefaultQuery("is_temp", "") != "" {
		tx.Where(fmt.Sprintf("%s.is_temp=?", model.TableDorPeople), query.IsTemp)
	}

	if query.Keyword != "" {
		if query.Type == "mobile" {
			tx.
				Joins(fmt.Sprintf("left join %s on %s.member_id=%s.id", model.TableMemMember, model.TableDorPeople, model.TableMemMember)).
				Where(fmt.Sprintf("%s.mobile", model.TableMemMember), query.Keyword)
		} else if query.Type == "room" {
			tx.
				Joins(fmt.Sprintf("left join %s on %s.room_id=%s.id", model.TableDorRoom, model.TableDorPeople, model.TableDorRoom)).
				Where(fmt.Sprintf("%s.name like ?", model.TableDorRoom), "%"+query.Keyword+"%")
		} else {
			tx.
				Joins(fmt.Sprintf("left join %s on %s.member_id=%s.id", model.TableMemMember, model.TableDorPeople, model.TableMemMember)).
				Where(fmt.Sprintf("%s.name", model.TableMemMember), query.Keyword)
		}
	}

	tc := tx

	responses := response.Paginate{
		Code:    20000,
		Message: "Success",
	}

	responses.Data.Size = query.GetSize()
	responses.Data.Page = query.GetPage()

	tc.Model(model.DorPeople{}).Count(&responses.Data.Total)

	if responses.Data.Total > 0 {
		var peoples []model.DorPeople

		tx.
			Preload("Member").
			Preload("Staff").
			Preload("Certification").
			Preload("Category").
			Preload("Building").
			Preload("Floor").
			Preload("Room").
			Preload("Bed").
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
				End:       item.End.ToDateString(),
				CreatedAt: item.CreatedAt.ToDateTimeString(),
			}
			if item.Staff.Id > 0 {
				results.Staff = item.Staff.Status
				results.Titles = item.Staff.Title
			}
			if item.Certification.Id > 0 {
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
	data.Database.Where("is_enable", constant.IsEnableYes).First(&bed)
	if bed.Id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "床位不存在",
		})
		return
	}

	var category model.DorStayCategory
	data.Database.Where("is_enable", constant.IsEnableYes).First(&category)
	if category.Id <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "类型不存在",
		})
		return
	}

	var count int64 = 0
	data.Database.Model(model.DorPeople{}).Joins(fmt.Sprintf("left join %s on %s.member_id=%s.id", model.TableMemMember, model.TableDorPeople, model.TableMemMember)).Where(fmt.Sprintf("%s.mobile", model.TableMemMember), former.Mobile).Where(fmt.Sprintf("%s.status", model.TableDorPeople), model.DorPeopleStatusLive).Count(&count)
	if count > 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    40400,
			Message: "该手机号已办理了入住，无法重复办理",
		})
		return
	}

	var member model.MemMember
	data.Database.Where("mobile", former.Mobile).First(&member)
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

	people := model.DorPeople{
		CategoryId: category.Id,
		BuildingId: bed.BuildingId,
		FloorId:    bed.FloorId,
		RoomId:     bed.RoomId,
		BedId:      bed.Id,
		TypeId:     bed.TypeId,
		MemberId:   member.Id,
		Start:      former.Start,
		End:        former.End,
		Status:     model.DorPeopleStatusLive,
		Remark:     former.Remark,
	}

	if t := data.Database.Create(&people); t.RowsAffected <= 0 {
		ctx.JSON(http.StatusOK, response.Response{
			Code:    60000,
			Message: "办理失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, response.Response{
		Code:    20000,
		Message: "Success",
	})
}
