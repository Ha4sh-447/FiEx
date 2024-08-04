package files

import (
	"io/fs"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"sync"

	"github.com/saracen/walker"
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

type DiskInfo struct {
	Mountpoint string
	Usage      *disk.UsageStat
}

func GetDiskUsage() ([]DiskInfo, error) {
	partitions, err := disk.Partitions(false)
	if err != nil {
		return nil, err
	}

	var diskInfoList []DiskInfo
	for _, partition := range partitions {
		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			continue // Ignore errors for individual partitions
		}
		diskInfoList = append(diskInfoList, DiskInfo{
			Mountpoint: partition.Mountpoint,
			Usage:      usage,
		})
	}
	return diskInfoList, nil
}

func TraverseDir(dir string) *sync.Map {
	// var fileBuff []string
	// var m map[string][]string
	// m := make(map[string][]string)
	m := &sync.Map{}

	dir, err := filepath.Abs(filepath.Clean(dir))
	if err != nil {
		slog.Error("ERROR", "ABSOLUTE ERROR", err)
	}
	m.Store(dir, []string{})

	// this can be optimized?

	walkFun := func(pathname string, fi os.FileInfo) error {
		pathname, err := filepath.Abs(filepath.Clean(pathname))
		if err != nil {
			slog.Error("ERROR", "Absolute file path error", err)
		}
		rootDir := filepath.Clean(filepath.Dir(pathname))

		// whatever the path may come put it in it's root dir
		if load, isPresent := m.LoadOrStore(rootDir, []string{pathname}); isPresent {
			fileList := load.([]string)
			if !contains(fileList, pathname) {
				m.Store(rootDir, append(fileList, pathname))
			}
		}

		// Add pathname to all its ancestor directories up to the root
		for currentDir := rootDir; strings.HasPrefix(currentDir, dir); currentDir = filepath.Dir(currentDir) {
			if load, isPresent := m.LoadOrStore(currentDir, []string{pathname}); isPresent {
				fileList := load.([]string)
				if !contains(fileList, pathname) {
					m.Store(currentDir, append(fileList, pathname))
				}
			}

			// Stop if we have reached the root directory
			if currentDir == dir {
				break
			}
		}

		return nil
	}

	errorCallbackOption := walker.WithErrorCallback(func(pathname string, err error) error {
		// ignore permissione errors
		if os.IsPermission(err) {
			return nil
		}

		return err
	})

	walker.Walk(dir, walkFun, errorCallbackOption)

	return m
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
