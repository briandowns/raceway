package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// NewDeploymentParams holds the decoded JSON from a POST to create
// a new deployment
type NewDeploymentParams struct {
	FromENV bool   `json:"from_env"`
	Name    string `json:"name"`
}

// ScenarioContent provides the contents for a given scenario
func ScenarioContent(scenPath, scenario string) map[string]interface{} {
	results := make(map[string]interface{})
	walkFn := func(path string, info os.FileInfo, err error) error {
		stat, err := os.Stat(path)
		if err != nil {
			return err
		}

		if stat.IsDir() && path != scenPath && !true {
			return filepath.SkipDir
		}

		if filepath.Base(path) == scenario {
			f, err := os.Open(path)
			if err != nil {
				log.Fatalln(err)
			}
			decoder := json.NewDecoder(f)
			err = decoder.Decode(&results)
			if err != nil {
				return err
			}
		}
		return nil
	}
	err := filepath.Walk(scenPath, walkFn)
	if err != nil {
		log.Println(err)
	}
	return results
}

// Scenarios searches for all files matching the configuration
func Scenarios(scenPath string, scenFormat string) map[string][]string {
	foundMap := make(map[string][]string)
	walkFn := func(path string, info os.FileInfo, err error) error {
		stat, err := os.Stat(path)
		if err != nil {
			return err
		}

		if stat.IsDir() && path != scenPath && !true {
			return filepath.SkipDir
		}

		if filepath.Ext(path) == "."+scenFormat {
			sep := strings.Split(path, "/")
			if _, ok := foundMap[sep[len(sep)-2]]; !ok {
				foundMap[sep[len(sep)-2]] = append(foundMap[sep[len(sep)-2]], sep[len(sep)-1])
			} else {
				foundMap[sep[len(sep)-2]] = append(foundMap[sep[len(sep)-2]], sep[len(sep)-1])
			}
		}
		return nil
	}

	err := filepath.Walk(scenPath, walkFn)
	if err != nil {
		log.Println(err)
		return nil
	}
	return foundMap
}
