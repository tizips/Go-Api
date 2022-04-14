package dormitory

import (
	"fmt"
	"github.com/golang-module/carbon/v2"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
	"saas/app/model"
	"saas/kernel/data"
)

func CrontabDayPeople(cron *cron.Cron) {

	if _, err := cron.AddFunc("0 2,3 * * *", handler); err != nil {
		fmt.Printf("启动失败：%s", err.Error())
		return
	}

}

func handler() {

	now := carbon.Now()

	var peoples []model.DorPeople

	data.Database.
		Preload("Master", func(db *gorm.DB) *gorm.DB { return db.Unscoped() }).
		Where("`start`<=? and `status`=?", now.Yesterday().ToDateString(), model.DorPeopleStatusLive).
		Where("not exists (?)", data.Database.
			Select("1").
			Table(model.TableDorDay).
			Where(fmt.Sprintf("%s.`id`=%s.`people_id`", model.TableDorPeople, model.TableDorDay)).
			Where(fmt.Sprintf("%s.`date`=?", model.TableDorDay), now.Yesterday().ToDateString()),
		).
		FindInBatches(&peoples, 50, func(tx *gorm.DB, batch int) error {

			if len(peoples) > 0 {
				var days = make([]model.DorDay, len(peoples))
				for index, item := range peoples {
					days[index] = model.DorDay{
						CategoryId:     item.CategoryId,
						TypeId:         item.TypeId,
						BuildingId:     item.BuildingId,
						FloorId:        item.FloorId,
						RoomId:         item.RoomId,
						BedId:          item.BedId,
						PeopleId:       item.Id,
						MemberId:       item.MemberId,
						MasterMemberId: item.Master.MemberId,
						MasterPeopleId: item.MasterId,
						Date:           carbon.Date{Carbon: now.Yesterday()},
					}
				}
				data.Database.Create(&days)
			}

			return nil
		})

}
