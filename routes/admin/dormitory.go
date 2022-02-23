package admin

import (
	"github.com/gin-gonic/gin"
	"saas/app/controller/admin/domitory/basic"
	"saas/app/controller/admin/domitory/stay"
)

func RouteDormitory(route *gin.RouterGroup) {

	dormitory := route.Group("dormitory")
	{
		basicGroup := dormitory.Group("basic")
		{
			buildings := basicGroup.Group("buildings")
			{
				buildings.GET("", basic.ToBuildingByList)
				buildings.PUT(":id", basic.DoBuildingByUpdate)
				buildings.DELETE(":id", basic.DoBuildingByDelete)
			}

			building := basicGroup.Group("building")
			{
				building.POST("", basic.DoBuildingByCreate)
				building.PUT("enable", basic.DoBuildingByEnable)
				building.GET("online", basic.ToBuildingByOnline)
			}

			typ := basicGroup.Group("types")
			{
				typ.GET("", basic.ToTypeByList)
				typ.PUT(":id", basic.DoTypeByUpdate)
				typ.DELETE(":id", basic.DoTypeByDelete)
			}

			Type := basicGroup.Group("type")
			{
				Type.POST("", basic.DoTypeByCreate)
				Type.PUT("enable", basic.DoTypeByEnable)
				Type.GET("online", basic.ToTypeByOnline)
			}

			floors := basicGroup.Group("floors")
			{
				floors.GET("", basic.ToFloorByList)
				floors.PUT(":id", basic.DoFloorByUpdate)
				floors.DELETE(":id", basic.DoFloorByDelete)
			}

			floor := basicGroup.Group("floor")
			{
				floor.POST("", basic.DoFloorByCreate)
				floor.PUT("enable", basic.DoFloorByEnable)
				floor.GET("online", basic.ToFloorByOnline)
			}

			rooms := basicGroup.Group("rooms")
			{
				rooms.GET("", basic.ToRoomByPaginate)
				rooms.PUT(":id", basic.DoRoomByUpdate)
				rooms.DELETE(":id", basic.DoRoomByDelete)
			}

			room := basicGroup.Group("room")
			{
				room.POST("", basic.DoRoomByCreate)
				room.PUT("enable", basic.DoRoomByEnable)
				room.PUT("furnish", basic.DoRoomByFurnish)
				room.GET("online", basic.ToRoomByOnline)
			}

			beds := basicGroup.Group("beds")
			{
				beds.GET("", basic.ToBedByPaginate)
				beds.PUT(":id", basic.DoBedByUpdate)
				beds.DELETE(":id", basic.DoBedByDelete)
			}

			bed := basicGroup.Group("bed")
			{
				bed.POST("", basic.DoBedByCreate)
				bed.PUT("enable", basic.DoBedByEnable)
				bed.GET("online", basic.ToBedByOnline)
			}
		}

		stayGroup := dormitory.Group("stay")
		{
			categories := stayGroup.Group("categories")
			{
				categories.GET("", stay.ToCategoryByList)
				categories.PUT(":id", stay.DoCategoryByUpdate)
				categories.DELETE(":id", stay.DoCategoryByDelete)
			}

			category := stayGroup.Group("category")
			{
				category.POST("", stay.DoCategoryByCreate)
				category.POST("enable", stay.DoCategoryByEnable)
				category.GET("online", stay.ToCategoryByOnline)
			}

			peoples := stayGroup.Group("peoples")
			{
				peoples.GET("", stay.ToPeopleByPaginate)
			}

			people := stayGroup.Group("people")
			{
				people.POST("", stay.DoPeopleByCreate)
			}
		}
	}

}
