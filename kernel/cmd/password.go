package cmd

import (
	"github.com/gookit/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

func Password(command *cobra.Command) {

	cmd := &cobra.Command{
		Use:   "password",
		Short: "生成密码",
		Run:   password,
	}

	command.AddCommand(cmd)
}

func password(c *cobra.Command, args []string) {

	prompt := promptui.Prompt{
		Label: "密码",
	}

	psd, _ := prompt.Run()

	psd = strings.TrimSpace(psd)

	if psd == "" {
		color.Errorf("密码不能为空")
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(psd), bcrypt.DefaultCost)

	color.Success.Println(string(hash))
}
