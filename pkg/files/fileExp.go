package internal

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"reflect"
)

// Gives list of files contained in the current directory along with
// user object and an error

func Files(cwd string) (*user.User, []fs.DirEntry, error) {
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
		return nil, nil, err
	}

	files, err := os.ReadDir(cwd)
	if err != nil {
		log.Fatal(err)
		return nil, nil, err
	}

	fmt.Println(reflect.TypeOf(files))
	fmt.Println(user.HomeDir)

	// for _, file := range files {
	// 	fmt.Println(file.Name())
	// }
	fmt.Println(user)
	return user, files, nil
}

// Executes command to open directories/files
func OpenFile(path string) {
	switch filepath.Ext(path) {
	case ".txt":
		exec.Command("notepad", path).Start()
	default:
		// Handle other file types or open with the default application
		exec.Command("cmd", "/C", "start", path).Start()
	}
}
