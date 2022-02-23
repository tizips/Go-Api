package system

import "fmt"

func Help() {
	fmt.Println("")
	fmt.Println("Saas")
	fmt.Println("")

	fmt.Println("Usage:")
	fmt.Println("\tapp [command]")
	fmt.Println("")

	fmt.Println("The commands are:")
	fmt.Println("\thelp\t\t帮助")
	fmt.Println("\tserver\t\t启动服务")
	fmt.Println("\tmigrate\t\t执行数据迁移")
	fmt.Println("\trollback\t执行数据迁移回滚")
	fmt.Println("\tpassword\t登陆密码生成 -password=password")
	fmt.Println("")
}
