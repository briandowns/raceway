package database

import (
	"fmt"

	"github.com/briandowns/aion/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/jmoiron/sqlx"
)

// Database holds db conf and a connection
type Database struct {
	Conf *config.Config
	Conn *gorm.DB
}

// NewDatabase creates a new Database object
func NewDatabase(conf *config.Config) (*Database, error) {
	d := &Database{
		Conf: conf,
	}
	if err := d.connect(); err != nil {
		return nil, err
	}
	return d, nil
}

// Connect will provide the caller with a db connection
func (d *Database) connect() error {
	db, err := gorm.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?timeout=%s&charset=utf8&parseTime=True&loc=Local",
			d.Conf.Database.DBUser, d.Conf.Database.DBPass, d.Conf.Database.DBHost, d.Conf.Database.DBPort, d.Conf.Database.DBName, "60s"))
	if err != nil {
		return err
	}
	db.LogMode(true)
	d.Conn = &db
	return nil
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

// GetTasks gets all tasks from the database
func (d *Database) GetTasks() []Task {
	var data []Task
	d.Conn.Find(&data)
	return data
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
