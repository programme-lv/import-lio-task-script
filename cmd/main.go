package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/programme-lv/import-lio-task-script/internal"
)

func main() {
	// Define flags
	sourceDir := flag.String("source", "", "Source directory containing the tasks")
	sourceFormat := flag.String("format", "lio2024", "Source format of the tasks")
	destDir := flag.String("dest", "", "Destination directory where the new directory will be placed")

	// Parse flags
	flag.Parse()

	// Validate flags
	if *sourceDir == "" || *destDir == "" {
		fmt.Println("Source and destination directories must be specified.")
		flag.Usage()
		os.Exit(1)
	}

	// Check if the source format is "lio2024"
	if *sourceFormat != "lio2024" {
		fmt.Println("Unsupported source format. Only 'lio2024' is supported.")
		os.Exit(1)
	}

	// Get the base name of the source directory
	baseName := filepath.Base(*sourceDir)
	newDirName := baseName + "_proglv"
	newDirPath := filepath.Join(*destDir, newDirName)

	task, err := internal.ParseLio2024TaskDir(*sourceDir)
	if err != nil {
		log.Fatalf("Failed to parse Lio2024 task: %v\n", err)
	}

	err = task.Store(newDirPath)
	if err != nil {
		log.Fatalf("Failed to store task: %v\n", err)
	}
}
