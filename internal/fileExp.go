package internal

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/user"
	"reflect"
)

func Files(cwd string) (*user.User, []fs.DirEntry, error) {
	user, err := user.Current()
	files, err := os.ReadDir(cwd)

	fmt.Println(reflect.TypeOf(files))
	fmt.Println(user.HomeDir)

	if err != nil {
		log.Fatal(err)
		return nil, nil, err
	}

	// for _, file := range files {
	// 	fmt.Println(file.Name())
	// }
	fmt.Println(user)
	return user, files, nil
}
