package main

import (
	"bufio"
	"os/exec"
	"regexp"
	"strings"
)

const CMD = "rally"

var r = regexp.MustCompile("[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}")

// Run exeutes the given task
func TaskStart(scenario string, jsonOutput, HTMLOutput bool, cmdStartArgs ...string) error {
	cmdArgs := []string{"task", "start", scenario}
	cmd := exec.Command(CMD, cmdStartArgs...)
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	var taskID string
	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			if r.MatchString(scanner.Text()) {
				taskID = strings.Trim(strings.Split(scanner.Text(), " ")[3], ":")
				return
			}
		}
	}()
	if err = cmd.Start(); err != nil {
		return err
	}
	if err = cmd.Wait(); err != nil {
		return err
	}
	if jsonOutput {
		if err = generateJSON(taskID); err != nil {
			return err
		}
	}
	if HTMLOutput {
		if err = generateHTML(taskID); err != nil {
			return err
		}
	}
	return nil
}

type Generator interface {
	Generate(output string) error
}

// generateHTML builds HTML pages based
func generateHTML(taskID string, cmdStartArgs ...string) error {
	cmdArgs := []string{"task", "report", taskID, "--output", Conf.HTMLOutputDir + "/" + taskID + ".html"}
	cmd := exec.Command(CMD, cmdStartArgs...)
	return nil
}

// generateJSON
func generateJSON(taskID string, cmdStartArgs ...string) error {
	cmdArgs := []string{"task", "results", taskID}
	cmd := exec.Command(CMD, cmdStartArgs...)
	return nil
}
