package files

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"reflect"

	"github.com/shirou/gopsutil/disk"
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

func Volume() {
	partitions, err := disk.Partitions(false)
	if err != nil {
		log.Fatalf("Failed to get disk partitions: %v", err)
	}

	for _, partition := range partitions {
		fmt.Printf("Listing files in %s:\n", partition.Mountpoint)
		files, err := os.ReadDir(partition.Mountpoint + "\\")
		if err != nil {
			log.Fatalf("Failed to read the directory %s: %v", partition.Mountpoint, err)
		}

		for _, file := range files {
			fmt.Println(file.Name())
		}
	}
}
