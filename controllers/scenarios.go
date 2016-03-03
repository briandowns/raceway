package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"

	"github.com/briandowns/raceway/config"
)

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

// ScenariosHandler
func ScenariosHandler(ren *render.Render, conf *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ren.JSON(w, http.StatusOK, Scenarios(conf.Scenarios.ScenarioDir, conf.Scenarios.ScenarioFormat))
	}
}

// ScenariosByNameHandler
func ScenariosByNameHandler(ren *render.Render, conf *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		scenarioName := vars["name"]
		ren.JSON(w, http.StatusOK, ScenarioContent(conf.Scenarios.ScenarioDir, scenarioName))
	}
}
