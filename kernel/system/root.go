package system

import (
	"flag"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"os"
	"saas/app/constant"
	"saas/app/model"
	"saas/kernel/auth"
	"saas/kernel/data"
)

func Root() {

	var username, mobile, nickname, password string

	flag.StringVar(&username, "username", "", "Username")
	flag.StringVar(&mobile, "mobile", "", "Username")
	flag.StringVar(&nickname, "nickname", "", "Username")
	flag.StringVar(&password, "password", "", "")

	_ = flag.CommandLine.Parse(os.Args[2:])

	var SysAdmin model.SysAdmin

	if username == "" {
		fmt.Println("Username 不能为空")
		return
	} else {
		data.Database.Where("username", username).First(&SysAdmin)
	}

	if SysAdmin.Id <= 0 {

		if mobile == "" {
			fmt.Println("Mobile 不能为空")
			return
		}
		if nickname == "" {
			fmt.Println("Nickname 不能为空")
			return
		}
		if password == "" {
			fmt.Println("Password 不能为空")
			return
		}

		hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		SysAdmin = model.SysAdmin{
			Username: username,
			Mobile:   mobile,
			Nickname: nickname,
			Password: string(hash),
			IsEnable: constant.IsEnableYes,
		}

		tx := data.Database.Create(&SysAdmin)
		if tx.RowsAffected <= 0 {
			fmt.Printf("Create SysAdmin error:%v", tx.Error)
			return
		}
	}

	var SysAdminBindRole = model.SysAdminBindRole{AdminId: SysAdmin.Id, RoleId: auth.ROOT}
	data.Database.Where("admin_id", SysAdmin.Id).Where("role_id", auth.ROOT).FirstOrCreate(&SysAdminBindRole)
	if SysAdminBindRole.Id <= 0 {
		fmt.Println("Create admin bind role fail")
		return
	}

	exist, err := auth.Casbin.HasRoleForUser(auth.NameByAdmin(SysAdmin.Id), auth.NameByRole(auth.ROOT))
	if exist {
		fmt.Printf("Root of %s has existed", SysAdmin.Nickname)
		return
	} else if exist && err != nil {
		fmt.Printf("Root read error:%v", err)
		return
	}

	if ok, err := auth.Casbin.AddRoleForUser(auth.NameByAdmin(SysAdmin.Id), auth.NameByRole(auth.ROOT)); !ok {
		fmt.Printf("Root create fail: %v", err.Error())
		return
	}

	fmt.Println("Create success:", SysAdmin.Nickname)
}
