package system

import (
	"github.com/gookit/color"
	"github.com/manifoldco/promptui"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

func Password() {

	prompt := promptui.Prompt{
		Label: "密码",
	}

	password, _ := prompt.Run()

	password = strings.TrimSpace(password)

	if password == "" {
		color.Error.Println("密码不能为空")
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	color.Success.Println(string(hash))
}
