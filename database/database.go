package database

import (
	"fmt"

	"github.com/briandowns/raceway/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
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
	db.LogMode(d.Conf.Database.DBDebug)
	d.Conn = &db
	return nil
}

// GetDeployments gets all deployments Rally is aware of
func (d *Database) GetDeployments() []Deployment {
	var data []Deployment
	d.Conn.Find(&data)
	return data
}

// DeploymentsByName gets all deployments Rally is aware of
func (d *Database) DeploymentsByName() []Deployment {
	var data []Deployment
	return d.Conn.Find(&data)
	//return data
}

// TaskByUUID gets a task by its UUID
func (d *Database) TaskByUUID(uuid string) []Task {
	var data []Task
	d.Conn.Where("uuid = ?", uuid).Find(&data)
	return data
}

// GetTasks gets all tasks from the database
func (d *Database) GetTasks() []Task {
	var data []Task
	d.Conn.Find(&data)
	return data
}

// TasksRunning gets all tasks Rally is aware of with a status of 'running'
func (d *Database) TasksRunning() []Task {
	var data []Task
	return data
}

// TaskResultsByUUID gets the results of a task by ID after it's been run
func (d *Database) TaskResultsByUUID() []TaskResult {
	var data []TaskResult
	return data
}
