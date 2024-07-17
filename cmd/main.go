package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/programme-lv/import-lio-task-script/internal"

	"github.com/programme-lv/fs-task-problem-toml/pkg/ptoml"
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
	err = internal.Unzip(zipPath, filepath.Join(newDirPath, "tests"))
	if err != nil {
		log.Fatalf("Failed to unzip %s: %v\n", zipPath, err)
	}

	// Move PDF files
	pdfSourceDir := filepath.Join(*sourceDir, "teksts")
	pdfStatementDir := filepath.Join(newDirPath, "statements", "pdf")
	pdfDestPath := filepath.Join(pdfStatementDir, "lv.pdf")
	err = internal.CopyPDF(pdfSourceDir, pdfDestPath)
	if err != nil {
		log.Fatalf("Failed to move PDF files: %v\n", err)
	}

	// Read source task.yaml file
	taskYAMLPath := filepath.Join(*sourceDir, "task.yaml")
	task, err := internal.ReadTaskYAML(taskYAMLPath)
	if err != nil {
		log.Fatalf("Failed to read task.yaml: %v\n", err)
	}

	olympiad := "LIO"

	// Write problem.toml file
	problemToml := ptoml.ProblemTOMLV2dot2{
		TaskName:    task.Title,
		Metadata:    ptoml.ProblemTOMLV2dot0Metadata{ProblemTags: []string{}, DifficultyFrom1To5: 0, TaskAuthors: []string{}, OriginOlympiad: &olympiad},
		Constraints: ptoml.ProblemTOMLV2dot0Constraints{MemoryMegabytes: task.MemoryLimit, CPUTimeSeconds: task.TimeLimit},
		TestGroups:  []ptoml.ProblemTOMLV2dot1LIOTestGroup{},
		VisInpSTs:   []int{1},
	}

	// Read all filenames in the tests directory
	mapGroupToTestFilenames := map[int][]string{}
	testsDir := filepath.Join(newDirPath, "tests")
	files, err := os.ReadDir(testsDir)
	if err != nil {
		log.Fatalf("Failed to read directory %s: %v\n", testsDir, err)
	}
	for _, file := range files {
		if file.IsDir() {
			log.Fatalf("Unexpected directory in the tests directory: %s\n", file.Name())
		}
		fname := file.Name()
		// split into part by dot. keep the last part (the extension)
		parts := strings.Split(fname, ".")
		if len(parts) < 2 {
			log.Fatalf("Unexpected filename: %s\n", fname)
		}

		ext := parts[len(parts)-1]

		re := regexp.MustCompile("[0-9]+")
		dGroups := re.FindAllString(ext, -1)
		if len(dGroups) != 1 {
			log.Fatalf("Unexpected filename: %s\n", fname)
		}

		group, err := strconv.Atoi(dGroups[0])
		if err != nil {
			log.Fatalf("Failed to convert %s to int: %v\n", dGroups[0], err)
		}

		// rename the file. if it has ".i" in the name, append "in" to the end
		// if it has ".o" in the name, append "out" to the end
		// remove ".i" and ".o" from the name

		if strings.Contains(fname, ".i") {
			newFname := strings.ReplaceAll(fname, ".i", "")
			newFname += ".in"
			err := os.Rename(filepath.Join(testsDir, fname), filepath.Join(testsDir, newFname))
			if err != nil {
				log.Fatalf("Failed to rename %s to %s: %v\n", fname, newFname, err)
			}
			fname = newFname
		} else if strings.Contains(fname, ".o") {
			newFname := strings.ReplaceAll(fname, ".o", "")
			newFname += ".out"
			err := os.Rename(filepath.Join(testsDir, fname), filepath.Join(testsDir, newFname))
			if err != nil {
				log.Fatalf("Failed to rename %s to %s: %v\n", fname, newFname, err)
			}
			fname = newFname
		}

		mapGroupToTestFilenames[group] = append(mapGroupToTestFilenames[group], fname)
	}

	for _, group := range task.TestsGroups {
		groups := []int{}

		switch v := group.Groups.(type) {
		case int:
			groups = append(groups, v)
		case []interface{}:
			integers := []int{}
			for _, vv := range v {
				switch vv := vv.(type) {
				case int:
					integers = append(integers, vv)
				default:
					log.Fatalf("Unsupported group: %+v %T\n", vv, vv)
				}
			}
			if len(integers) == 1 {
				groups = append(groups, integers...)
			} else if len(v) == 2 {
				for i := integers[0]; i <= integers[1]; i++ {
					groups = append(groups, i)
				}
			} else {
				log.Fatalf("Unsupported groups length: %v\n", v)
			}
		default:
			log.Fatalf("Unsupported groups: %+v %T\n", v, v)
		}

		publicGroups := []int{}
		switch v := group.Public.(type) {
		case bool:
			if v {
				publicGroups = append(publicGroups, groups...)
			}
		case []interface{}:
			integers := []int{}
			for _, vv := range v {
				switch vv := vv.(type) {
				case int:
					integers = append(integers, vv)
				default:
					log.Fatalf("Unsupported public group: %+v %T\n", vv, vv)
				}
			}

			if len(integers) == 1 {
				publicGroups = append(publicGroups, integers...)
			} else if len(integers) == 2 {
				for i := integers[0]; i <= integers[1]; i++ {
					publicGroups = append(publicGroups, i)
				}
			} else {
				log.Fatalf("Unsupported public groups: %v\n", v)
			}
		default:
			log.Fatalf("Unsupported public groups: %v\n", v)
		}

		for _, g := range groups {
			// move g 0 files to examples
			if g == 0 {
				for _, f := range mapGroupToTestFilenames[g] {
					err := os.MkdirAll(filepath.Join(newDirPath, "examples"), 0755)
					if err != nil {
						log.Fatalf("Failed to create examples directory: %v\n", err)
					}
					err = os.Rename(filepath.Join(testsDir, f), filepath.Join(newDirPath, "examples", f))
					if err != nil {
						log.Fatalf("Failed to move %s to examples: %v\n", f, err)
					}
				}
				continue
			}
			// copy g 1 INPUT files to examples (outputs are hidden)
			// if g == 1 {
			// 	for _, f := range mapGroupToTestFilenames[g] {
			// 		if strings.Contains(f, ".o") {
			// 			continue
			// 		}
			// 		err := os.MkdirAll(filepath.Join(newDirPath, "examples"), 0755)
			// 		if err != nil {
			// 			log.Fatalf("Failed to create examples directory: %v\n", err)
			// 		}
			// 		content, err := os.ReadFile(filepath.Join(testsDir, f))
			// 		if err != nil {
			// 			log.Fatalf("Failed to read %s: %v\n", f, err)
			// 		}
			// 		err = os.WriteFile(filepath.Join(newDirPath, "examples", f), content, 0644)
			// 		if err != nil {
			// 			log.Fatalf("Failed to copy %s to examples: %v\n", f, err)
			// 		}
			// 	}
			// }

			isPublic := false
			for _, pg := range publicGroups {
				if g == pg {
					isPublic = true
					break
				}
			}
			problemToml.TestGroups = append(problemToml.TestGroups, ptoml.ProblemTOMLV2dot1LIOTestGroup{
				GroupID:    g,
				Points:     group.Points,
				Subtask:    group.Subtask,
				Public:     isPublic,
				TestFnames: mapGroupToTestFilenames[g],
			})
		}
	}

	output, err := problemToml.Marshall()
	if err != nil {
		log.Fatalf("Failed to marshal the problem.toml: %v\n", err)
	}

	err = os.WriteFile(filepath.Join(newDirPath, "problem.toml"), output, 0644)
	if err != nil {
		log.Fatalf("Failed to write problem.toml: %v\n", err)
	}
	log.Println("All done!")
}
