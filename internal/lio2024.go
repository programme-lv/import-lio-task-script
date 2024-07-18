package internal

import (
	"fmt"
	"os"
	"path/filepath"

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

	task, err := fstaskparser.NewTask(parsedYaml.FullTaskName)
	if err != nil {
		return nil, fmt.Errorf("failed to create new task: %v", err)
	}

	testZipAbsolutePath := filepath.Join(taskYamlPath, parsedYaml.TestZipPathRelToYaml)

	tests, err := ReadLioTestsFromZip(testZipAbsolutePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read tests from zip: %v", err)
	}

	for _, t := range tests {
		id := task.AddTest(t.Input, t.Answer)
		name := fmt.Sprintf("%03d_%d", t.TestGroup, t.NoInTestGroup)
		task.AssignFilenameToTest(name, id)
	}

	return task, nil
}
