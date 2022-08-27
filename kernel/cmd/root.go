package cmd

import (
	"errors"
	"github.com/gookit/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"saas/app/constant"
	"saas/app/model"
	"saas/kernel/app"
	"saas/kernel/authorize"
)

func Root(command *cobra.Command) {

	cmd := &cobra.Command{
		Use:   "root",
		Short: "生成系统开发账号",
		Run:   root,
	}

	command.AddCommand(cmd)

}

func root(cmd *cobra.Command, args []string) {

	var username, mobile, nickname, password string

	prompt := promptui.Prompt{
		Label: "账号",
		Validate: func(input string) error {

			if ok, _ := regexp.MatchString(`^[a-zA-Z\d\-_]{4,20}$`, input); !ok {
				return errors.New("请输入 4-20 位的英文数字以及 -_ 等字符")
			}

			return nil
		},
	}
	if username, _ = prompt.Run(); username == "" {
		return
	}

	var admin model.SysAdmin

	app.MySQL.Where("username=?", username).Find(&admin)

	if admin.Id <= 0 {

		prompt = promptui.Prompt{
			Label: "手机号",
			Validate: func(input string) error {

				if ok, _ := regexp.MatchString(`^1\d{10}$`, input); !ok {
					return errors.New("格式错误")
				}

				var total int64
				if app.MySQL.Model(&model.SysAdmin{}).Where("mobile=?", input).Count(&total); total > 0 {
					return errors.New("手机号已被使用")
				}

				return nil
			},
		}
		if mobile, _ = prompt.Run(); mobile == "" {
			return
		}

		prompt = promptui.Prompt{
			Label: "密码",
			Validate: func(input string) error {

				if ok, _ := regexp.MatchString(`^[a-zA-Z\d\-_@$&%!]{6,32}$`, input); !ok {
					return errors.New("请输入 6-32 位的英文数字以及 -_@$&%! 等字符")
				}

				return nil
			},
		}
		if password, _ = prompt.Run(); password == "" {
			return
		}

		prompt = promptui.Prompt{
			Label: "昵称",
			Validate: func(input string) error {

				if len(input) < 2 || len(input) > 32 {
					return errors.New("最多输入 2-32 个字符")
				}

				return nil
			},
		}
		if nickname, _ = prompt.Run(); nickname == "" {
			return
		}

		hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		admin = model.SysAdmin{
			Username: username,
			Mobile:   mobile,
			Nickname: nickname,
			Password: string(hash),
			IsEnable: constant.IsEnableYes,
		}

		if tx := app.MySQL.Create(&admin); tx.RowsAffected <= 0 {
			color.Errorf("Create Admin error: %v", tx.Error)
			return
		}
	}

	if ok, _ := app.Casbin.HasRoleForUser(authorize.NameByAdmin(admin.Id), authorize.NameByRole(authorize.ROOT)); ok {
		color.Errorf("账号：%s 已经有 ROOT 权限", username)
		return
	}

	var bind = model.SysAdminBindRole{AdminId: admin.Id, RoleId: authorize.ROOT}
	app.MySQL.Where("`admin_id`=? and `role_id`=?", admin.Id, authorize.ROOT).FirstOrCreate(&bind)
	if bind.Id <= 0 {
		color.Errorln("Create Admin bind Role fail")
		return
	}

	if ok, err := app.Casbin.AddRoleForUser(authorize.NameByAdmin(admin.Id), authorize.NameByRole(authorize.ROOT)); !ok {
		color.Errorf("Root create fail: %v", err)
		return
	}

	color.Success.Printf("Create success: %s", admin.Nickname)
}
