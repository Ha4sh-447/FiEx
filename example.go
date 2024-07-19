package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/charlievieth/fastwalk"
	"github.com/saracen/walker"
)

func goWalkThroughFiles(dir string) ([]string, error) {
	var fileBuff []string
	var mu sync.Mutex
	var wg sync.WaitGroup

	start := time.Now()

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if os.IsPermission(err) {
				// Silently skip permission errors
				return fs.SkipDir
			}
			return err
		}

		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			mu.Lock()
			fileBuff = append(fileBuff, p)
			mu.Unlock()
		}(path)

		return nil
	})
	if err != nil {
		log.Fatal("Error getting files: ", err)
	}

	wg.Wait()

	// fmt.Println(fileBuff)
	fmt.Println("go walk through files: ", time.Since(start))
	return fileBuff, nil
}

func walkerReadFiles(path string) ([]string, error) {
	start := time.Now()
	var fileBuff []string
	// Read callback
	walkFun := func(pathname string, fi os.FileInfo) error {
		fileBuff = append(fileBuff, pathname)
		return nil
	}

	// error callback option
	errorCallbackOption := walker.WithErrorCallback(func(pathname string, err error) error {
		// ignore permissione errors
		if os.IsPermission(err) {
			return nil
		}
		// halt traversal on any other error
		return err
	})

	cpuLimit := walker.WithLimit(16)

	walker.Walk(path, walkFun, errorCallbackOption, cpuLimit)
	fmt.Println("Time Taken by Walker function: ", time.Since(start))
	return fileBuff, nil
}

func walkDirFastWalk(dir string) ([]string, error) {
	start := time.Now()
	var fileBuff []string
	// Read callback
	conf := fastwalk.Config{
		NumWorkers: 100, // Use a more reasonable number of workers
	}

	walkFun := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", path, err)
			if os.IsPermission(err) {
				return fastwalk.SkipDir
			}
			return nil // returning nil allows the walk to continue
		}
		fileBuff = append(fileBuff, path)
		return nil
	}

	err := fastwalk.Walk(&conf, dir, walkFun)
	if err != nil {
		return nil, err
	}

	fmt.Println("Time taken by fast walk: ", time.Since(start))
	return fileBuff, nil
}

// READ JSON

func ReadJSONFile(filePath string) (map[string][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data map[string][]string
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}

	return data, nil
}

func main() {
	// Specify the path to the JSON file
	filePath := "C:\\Users\\harsh\\AppData\\Local\\__recent_cache__.json"

	// Read the JSON file
	data, err := ReadJSONFile(filePath)
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}

	fmt.Println(data)
	fmt.Println("----------------------------")

	fmt.Println(data["main.go"])
	// Print the data
	// for key, paths := range data {
	// 	fmt.Println("File:", key)
	// 	for _, path := range paths {
	// 		fmt.Println(" -", path)
	// 	}
	// }
}
