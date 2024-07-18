package internal

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/programme-lv/fs-task-format-parser/pkg/fstaskparser"
)

func ParseLio2024TaskDir(dirPath string) (*fstaskparser.Task, error) {
	taskYamlPath := filepath.Join(dirPath, "task.yaml")

	taskYamlContent, err := os.ReadFile(taskYamlPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read task.yaml: %v", err)
	}

	parsedYaml, err := ParseLio2024Yaml(taskYamlContent)
	if err != nil {
		return nil, fmt.Errorf("failed to parse task.yaml: %v", err)
	}

	if parsedYaml.CheckerPathRelToYaml != nil {
		// TODO: implement
		log.Fatalf("found checker %s", *parsedYaml.CheckerPathRelToYaml)
		return nil, fmt.Errorf("checkers are not implemented yet")
	}

	if parsedYaml.InteractorPathRelToYaml != nil {
		// TODO: implement
		return nil, fmt.Errorf("interactors are not implemented yet")
	}

	task, err := fstaskparser.NewTask(parsedYaml.FullTaskName)
	if err != nil {
		return nil, fmt.Errorf("failed to create new task: %v", err)
	}

	testZipAbsolutePath := filepath.Join(dirPath, parsedYaml.TestZipPathRelToYaml)

	tests, err := ReadLioTestsFromZip(testZipAbsolutePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read tests from zip: %v", err)
	}

	sort.Slice(tests, func(i, j int) bool {
		if tests[i].TestGroup == tests[j].TestGroup {
			return tests[i].NoInTestGroup < tests[j].NoInTestGroup
		}
		return tests[i].TestGroup < tests[j].TestGroup
	})

	mapTestsToTestGroups := map[int][]int{}

	for _, t := range tests {
		if t.TestGroup == 0 {
			task.AddExample(t.Input, t.Answer)
			continue
		}
		id := task.AddTest(t.Input, t.Answer)
		name := fmt.Sprintf("%03d_%d", t.TestGroup, t.NoInTestGroup)
		task.AssignFilenameToTest(name, id)
		mapTestsToTestGroups[t.TestGroup] = append(mapTestsToTestGroups[t.TestGroup], id)
	}

	for _, g := range parsedYaml.TestGroups {
		if g.GroupID == 0 {
			continue // examples
		}
		err := task.AddTestGroupWithID(g.GroupID, g.Points,
			g.Public, mapTestsToTestGroups[g.GroupID],
			g.Subtask)
		if err != nil {
			return nil, fmt.Errorf("failed to add test group: %v", err)
		}
	}

	/*
			type ParsedLio2024Yaml struct {
			CheckerPathRelToYaml    string
			InteractorPathRelToYaml string
		}
	*/
	task.SetCPUTimeLimitInSeconds(parsedYaml.CpuTimeLimitInSeconds)
	task.SetMemoryLimitInMegabytes(parsedYaml.MemoryLimitInMegabytes)

	pdfFilePath := filepath.Join(dirPath, "teksts")
	pdfFiles, err := filepath.Glob(filepath.Join(pdfFilePath, "*.pdf"))
	if err != nil {
		return nil, fmt.Errorf("failed to find PDF files: %w", err)
	}

	if len(pdfFiles) == 0 {
		return nil, fmt.Errorf("no PDF files found in the directory %s", pdfFilePath)
	}

	if len(pdfFiles) > 1 {
		return nil, fmt.Errorf("more than one PDF file found in the directory (%d)", len(pdfFiles))
	}

	pdfStatementPath := pdfFiles[0]

	pdfBytes, err := os.ReadFile(pdfStatementPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read PDF file: %w", err)
	}

	err = task.AddPDFStatement("lv", pdfBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to add PDF statement: %w", err)
	}

	task.AddVisibleInputSubtask(1)
	task.SetOriginOlympiad("LIO")

	// TODO: implement adding checker and interactor if present

	return task, nil
}
