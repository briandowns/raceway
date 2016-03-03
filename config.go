package main

import (
	"encoding/json"
	"os"
)

// Configuration contains the Raceway configuration
type Configuration struct {
	Raceway struct {
		AppPort  string `json:"app_port"`
		AppDebug bool   `json:"app_debug"`
	} `json:"raceway"`
	MySQL struct {
		DBHost string `json:"db_host"`
		DBPort int    `json:"db_port"`
		DBUser string `json:"db_user"`
		DBPass string `json:"db_pass"`
		DBName string `json:"db_name"`
	} `json:"mysql"`
	Scenarios struct {
		ScenarioDir    string `json:"scenario_dir"`
		ScenarioFormat string `json:"scenario_format"`
	} `json:"scenarios"`
	Graphite struct {
		GraphiteHost string `json:"graphite_host"`
	} `json:"graphite"`
	Scheduler struct {
		SchedulerDBBucket  string `json:"scheduler_db_bucket"`
		SchedulerDBName    string `json:"scheduler_db_name"`
		SchedulerDBTimeout int    `json:"scheduler_db_timeout"`
	} `json:"scheduler"`
	HTMLOutputDir string `json:"html_output_dir"`
}

// Load builds a config obj
func Load(cf string) (*Configuration, error) {
	confFile, err := os.Open(cf)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(confFile)
	conf := &Configuration{}
	err = decoder.Decode(conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}
