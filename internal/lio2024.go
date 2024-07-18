package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/programme-lv/fs-task-format-parser/pkg/fstaskparser"
)

func ReadLio2024TaskDir(dirPath string) (*fstaskparser.Task, error) {
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

	testZipAbsolutePath := filepath.Join(taskYamlPath, parsedYaml.TestZipPathRelToYaml)

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
		id := task.AddTest(t.Input, t.Answer)
		name := fmt.Sprintf("%03d_%d", t.TestGroup, t.NoInTestGroup)
		task.AssignFilenameToTest(name, id)
		mapTestsToTestGroups[t.TestGroup] = append(mapTestsToTestGroups[t.TestGroup], id)
	}

	for _, g := range parsedYaml.TestGroups {
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

	// TODO: implement adding pdf statement, checker and interactor if present

	pdfName := fmt.Sprintf("%s.pdf", parsedYaml.TaskShortIDCode)
	pdfStatementPath := filepath.Join(dirPath, "teksts", pdfName)

	pdfBytes, err := os.ReadFile(pdfStatementPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read PDF statement: %v", err)
	}

	err = task.AddPDFStatement("lv", pdfBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to add PDF statement: %v", err)
	}

	// TODO: implement adding checker and interactor if present

	return task, nil
}
