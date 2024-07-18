package internal

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

type ParsedLio2024Yaml struct {
	CpuTimeLimitInSeconds   float64
	MemoryLimitInMegabytes  int
	FullTaskName            string
	TestZipPathRelToYaml    string
	CheckerPathRelToYaml    string
	InteractorPathRelToYaml string
	SubtaskPoints           []int
	TestGroups              []ParsedLio2024YamlTestGroup
}

type ParsedLio2024YamlTestGroup struct {
	GroupID int
	Points  int
	Public  bool
	Subtask int
	Comment *string
}

type lio2024RawYaml struct {
	TimeLimit         float64                   `yaml:"time_limit"`
	MemoryLimit       int                       `yaml:"memory_limit"`
	ShortCode         string                    `yaml:"name"`
	TaskName          string                    `yaml:"title"`
	TestsZipRelPath   string                    `yaml:"tests_archive"`
	CheckerRelPath    string                    `yaml:"checker"`
	InteractorRelPath string                    `yaml:"interactor"`
	SubtaskPoitns     []int                     `yaml:"subtask_points"`
	TestGroups        []lio2024RawYamlTestGroup `yaml:"tests_groups"`
}

type lio2024RawYamlTestGroup struct {
	Groups  interface{} `yaml:"groups"`
	Points  int         `yaml:"points"`
	Public  interface{} `yaml:"public"`
	Subtask int         `yaml:"subtask"`
	Comment string      `yaml:"comment,omitempty"`
}

func ParseLio2024Yaml(content []byte) (res ParsedLio2024Yaml, err error) {
	rawYaml := lio2024RawYaml{}

	err = yaml.Unmarshal(content, &rawYaml)
	if err != nil {
		return
	}

	res.FullTaskName = rawYaml.TaskName
	res.CpuTimeLimitInSeconds = rawYaml.TimeLimit
	res.MemoryLimitInMegabytes = rawYaml.MemoryLimit
	res.TestZipPathRelToYaml = rawYaml.TestsZipRelPath
	res.CheckerPathRelToYaml = rawYaml.CheckerRelPath
	res.InteractorPathRelToYaml = rawYaml.InteractorRelPath
	res.SubtaskPoints = rawYaml.SubtaskPoitns

	for _, group := range rawYaml.TestGroups {
		groups := []ParsedLio2024YamlTestGroup{}

		allPublic := false
		isGroupPublic := make(map[int]bool)
		switch v := group.Public.(type) {
		case bool:
			allPublic = v
		case int:
			isGroupPublic[v] = true
		case []interface{}:
			integers := []int{}
			for _, vv := range v {
				switch vv := vv.(type) {
				case int:
					integers = append(integers, vv)
				default:
					err = fmt.Errorf("unsupported public groups: %+v %T", vv, vv)
					return
				}
			}
			if len(integers) != 1 {
				err = fmt.Errorf("unsupported public groups length: %v", v)
				return
			}
			for _, vv := range integers {
				isGroupPublic[vv] = true
			}
		}

		switch v := group.Groups.(type) {
		case int:
			var comment *string = nil
			if group.Comment != "" {
				comment = &group.Comment
			}
			groups = append(groups, ParsedLio2024YamlTestGroup{
				GroupID: v,
				Points:  group.Points,
				Public:  isGroupPublic[v] || allPublic,
				Subtask: group.Subtask,
				Comment: comment,
			})
		case []interface{}:
			integers := []int{}
			for _, vv := range v {
				switch vv := vv.(type) {
				case int:
					integers = append(integers, vv)
				default:
					err = fmt.Errorf("unsupported groups: %+v %T", vv, vv)
					return
				}
			}
			if len(integers) == 1 {
				var comment *string = nil
				if group.Comment != "" {
					comment = &group.Comment
				}
				groups = append(groups, ParsedLio2024YamlTestGroup{
					GroupID: integers[0],
					Points:  group.Points,
					Public:  false || allPublic,
					Subtask: group.Subtask,
					Comment: comment,
				})
			} else if len(v) == 2 {
				for i := integers[0]; i <= integers[1]; i++ {
					var comment *string = nil
					if group.Comment != "" {
						comment = &group.Comment
					}
					groups = append(groups, ParsedLio2024YamlTestGroup{
						GroupID: i,
						Points:  group.Points,
						Public:  isGroupPublic[i] || allPublic,
						Subtask: group.Subtask,
						Comment: comment,
					})
				}
			} else {
				err = fmt.Errorf("unsupported groups length: %v", v)
				return
			}
		default:
			err = fmt.Errorf("unsupported groups: %+v %T", v, v)
			return
		}

		res.TestGroups = append(res.TestGroups, groups...)
	}

	return
}

/*
excerpt from previous code

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
		}
	}
*/
