package system

import (
	"flag"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func Password() {

	var password string

	flag.StringVar(&password, "password", "", "密码生成")

	flag.Parse()

	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	fmt.Println(string(hash))
}
