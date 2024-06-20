package main

import (
	"fmt"
	"log"
	"os"

	"github.com/shirou/gopsutil/disk"
)

func main() {
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
