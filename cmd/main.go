package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"log"
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

	// Create the new directory
	err := os.Mkdir(newDirPath, 0755)
	if err != nil {
		fmt.Printf("Failed to create directory %s: %v\n", newDirPath, err)
		os.Exit(1)
	}

	fmt.Printf("New directory created at %s\n", newDirPath)

	// Unzip tests.zip
	zipPath := filepath.Join(*sourceDir, "testi", "tests.zip")
	err = Unzip(zipPath, newDirPath)
	if err != nil {
		log.Fatalf("Failed to unzip %s: %v\n", zipPath, err)
	}

	// Move PDF files
	pdfSourceDir := filepath.Join(*sourceDir, "teksts")
	pdfDestPath := filepath.Join(newDirPath, "statements", "pdf", "lv.pdf")
	err = MovePDF(pdfSourceDir, pdfDestPath)
	if err != nil {
		log.Fatalf("Failed to move PDF files: %v\n", err)
	}
}
