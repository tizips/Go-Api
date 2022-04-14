package dir

import (
	"os"
	"saas/kernel/config"
)

func InitDir() {

	_ = Mkdir(config.Application.Public)

}

func Mkdir(path string) error {

	isMake := false
	stat, err := os.Stat(path)
	if err != nil {
		isMake = true
	}

	if !isMake && !stat.IsDir() {
		isMake = true
	}

	if isMake {
		err = os.Mkdir(path, 0750)
		return err
	}

	return nil
}

func Touch(filename string) (*os.File, error) {

	isMake := false
	stat, err := os.Stat(filename)
	if err != nil {
		isMake = true
	}

	if !isMake && stat.IsDir() {
		isMake = true
	}

	var file *os.File

	if isMake {
		file, err = os.Create(filename)
		return nil, err
	}

	return file, nil
}
