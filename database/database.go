package database

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// Connect will provide the caller with a db connection
func Connect(user string, pass string, host string, port int, database string) (*sqlx.DB, error) {
        db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, pass, host, port, database))
	if err != nil {
		return nil, err
	}
	return db, nil
}

// AllDeployments gets all deployments Rally is aware of
func AllDeployments(conn *sqlx.DB) ([]Deployments, error) {
	dep := []Deployments{}
	err := conn.Select(&dep, "SELECT created_at, updated_at, id, uuid, parent_uuid, name, started_at, completed_at, enum_deployments_status FROM deployments")
	if err != nil {
		return nil, err
	}
	return dep, nil
}

// DeploymentsByName gets all deployments Rally is aware of
func DeploymentsByName(conn *sqlx.DB, name string) ([]Deployments, error) {
	dep := []Deployments{}
	err := conn.Select(&dep, fmt.Sprintf("SELECT created_at, updated_at, id, uuid, parent_uuid, name, started_at, completed_at, enum_deployments_status FROM deployments where name = '%s'", name))
	if err != nil {
		return nil, err
	}
	return dep, nil
}

// AllTasks gets all tasks Rally is aware of
func AllTasks(conn *sqlx.DB) ([]Tasks, error) {
	tasks := []Tasks{}
	err := conn.Select(&tasks, "SELECT created_at, updated_at, id, uuid, status, verification_log, tag, deployment_uuid FROM tasks")
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

// TaskByUUID gets a task by its UUID
func TaskByUUID(conn *sqlx.DB, taskUUID string) ([]Tasks, error) {
	tasks := []Tasks{}
	err := conn.Select(&tasks, fmt.Sprintf("SELECT created_at, updated_at, id, uuid, status, verification_log, tag, deployment_uuid FROM tasks where uuid = '%s'", taskUUID))
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

// TasksRunning gets all tasks Rally is aware of with a status of 'running'
func TasksRunning(conn *sqlx.DB) ([]Tasks, error) {
	tasks := []Tasks{}
	err := conn.Select(&tasks, "SELECT created_at, updated_at, id, uuid, status, verification_log, tag, deployment_uuid FROM tasks where status = 'running'")
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

// TaskResultsByUUID gets the results of a task by ID after it's been run
func TaskResultsByUUID(conn *sqlx.DB, taskUUID string) ([]TaskResults, error) {
	results := []TaskResults{}
	err := conn.Select(&results, fmt.Sprintf("SELECT * FROM task_results where task_uuid = '%s'", taskUUID))
	if err != nil {
		return nil, err
	}
	return results, nil
}
