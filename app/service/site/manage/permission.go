package manage

import (
	"saas/app/model"
	"saas/app/response/admin/site/manage"
	"saas/kernel/app"
)

func TreePermission(module int, parent bool, simple bool) []manage.TreePermission {

	tx := app.MySQL

	if module > 0 {
		tx = tx.Where("`module_id` = ?", module)
	}

	if parent {
		tx = tx.Where("`parent_i2`<=? and `method`=? and `path`=?", 0, "", "")
	}

	var permissions []model.SysPermission

	tx.Find(&permissions)

	return HandlerTree(permissions, parent, simple)
}

func HandlerTree(permissions []model.SysPermission, parent bool, simple bool) []manage.TreePermission {

	responses := make([]manage.TreePermission, 0)

	if len(permissions) > 0 {

		var children1, children2 []model.SysPermission

		for _, item := range permissions {
			if item.ParentI2 > 0 {
				children2 = append(children2, item)
			} else if item.ParentI1 > 0 {
				children1 = append(children1, item)
			} else {
				temp := manage.TreePermission{
					Id:        item.Id,
					Name:      item.Name,
					Slug:      item.Slug,
					Method:    item.Method,
					Path:      item.Path,
					CreatedAt: item.CreatedAt.ToDateTimeString(),
				}

				if parent || simple {
					temp.Method = ""
					temp.Path = ""
					temp.CreatedAt = ""
				}

				responses = append(responses, temp)
			}
		}

		for index, item := range responses {
			for _, value := range children1 {
				if item.Id == value.ParentI1 {
					childrenI1 := manage.TreePermission{
						Id:        value.Id,
						Parents:   []int{value.ParentI1},
						Name:      value.Name,
						Slug:      value.Slug,
						Method:    value.Method,
						Path:      value.Path,
						CreatedAt: value.CreatedAt.ToDateTimeString(),
					}

					if !parent {
						for _, val := range children2 {
							if childrenI1.Id == val.ParentI2 {
								childrenI2 := manage.TreePermission{
									Id:        val.Id,
									Parents:   []int{val.ParentI1, val.ParentI2},
									Name:      val.Name,
									Slug:      val.Slug,
									Method:    val.Method,
									Path:      val.Path,
									CreatedAt: val.CreatedAt.ToDateTimeString(),
								}

								if simple {
									childrenI2.Parents = nil
									childrenI2.Method = ""
									childrenI2.Path = ""
									childrenI2.CreatedAt = ""
								}

								childrenI1.Children = append(childrenI1.Children, childrenI2)
							}
						}
					}

					if simple {
						childrenI1.Parents = nil
						childrenI1.Method = ""
						childrenI1.Path = ""
						childrenI1.CreatedAt = ""
					}

					responses[index].Children = append(responses[index].Children, childrenI1)
				}
			}
		}
	}

	return responses

}
