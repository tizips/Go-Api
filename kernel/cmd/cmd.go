package cmd

import "github.com/spf13/cobra"

func InitCmd() {

	command := &cobra.Command{
		Use:     "saas",
		Short:   "一站式应用解决方案",
		Version: "1.0.0",
	}

	Server(command)
	Root(command)
	Password(command)

	_ = command.Execute()
}
