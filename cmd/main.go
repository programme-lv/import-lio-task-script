package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
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

	// Get the base name of the source directory
	baseName := filepath.Base(*sourceDir)
	newDirName := baseName + "_proglv"
	newDirPath := filepath.Join(*destDir, newDirName)

	// Create the new directory
	err := os.Mkdir(newDirPath, 0755)
	if err != nil {
		fmt.Printf("Failed to create directory %s: %v\n", newDirPath, err)
		os.Exit(1)
	}

	fmt.Printf("New directory created at %s\n", newDirPath)
}
