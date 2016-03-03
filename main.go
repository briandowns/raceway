package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/pborman/uuid"

	"github.com/codegangsta/negroni"
	"github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"

	"github.com/briandowns/raceway/database"
)

var taskChan = make(chan string, 0)
var conf *Configuration

// Task holds task details once initialized
type Task struct{}

func main() {
	// Source in the configuration
	conf, err := Load("./config.json")
	if err != nil {
		log.Fatalln(err)
	}

	// Connect to the Rally MySQL database
	dbConn, err := database.Connect(conf.MySQL.DBUser, conf.MySQL.DBPass, conf.MySQL.DBHost, conf.MySQL.DBName)
	if err != nil {
		log.Fatalln(err)
	}
	defer dbConn.Close()

	// Start the scheduler
	StartScheduler(conf)

	ren := render.New(render.Options{Layout: "layout"})

	store := cookiestore.New([]byte(uuid.NewUUID().String()))
	n := negroni.New(
		negroni.NewRecovery(),
		negroni.NewLogger(),
		negroni.NewStatic(http.Dir("public")),
	)
	n.Use(sessions.Sessions("session", store))

	router := mux.NewRouter()

	// API routes
	router.HandleFunc("/api/v1/deployments", func(w http.ResponseWriter, r *http.Request) {
		d, err := database.AllDeployments(dbConn)
		if err != nil {
			log.Println(err)
		}
		ren.JSON(w, http.StatusOK, d)
	}).Methods("GET")

	router.HandleFunc("/api/v1/deployments/{name}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		depName := vars["name"]
		d, err := database.DeploymentsByName(dbConn, depName)
		if err != nil {
			log.Println(err)
		}
		ren.JSON(w, http.StatusOK, d)
	}).Methods("GET")

	router.HandleFunc("/api/v1/deployments", func(w http.ResponseWriter, r *http.Request) {
		dp := &NewDeploymentParams{}
		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(&dp)
		if err != nil {
			log.Println(err)
		}
		ren.JSON(w, http.StatusOK, dp)
	}).Methods("POST")

	router.HandleFunc("/api/v1/scenarios", func(w http.ResponseWriter, r *http.Request) {
		ren.JSON(w, http.StatusOK, Scenarios(conf.Scenarios.ScenarioDir, conf.Scenarios.ScenarioFormat))
	}).Methods("GET")

	router.HandleFunc("/api/v1/scenarios/{name}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		scenarioName := vars["name"]
		ren.JSON(w, http.StatusOK, ScenarioContent(conf.Scenarios.ScenarioDir, scenarioName))
	}).Methods("GET")

	router.HandleFunc("/api/v1/tasks", func(w http.ResponseWriter, r *http.Request) {
		t, err := database.AllTasks(dbConn)
		if err != nil {
			log.Println(err)
		}
		ren.JSON(w, http.StatusOK, t)
	}).Methods("GET")

	router.HandleFunc("/api/v1/tasks/start", func(w http.ResponseWriter, r *http.Request) {
		ren.JSON(w, http.StatusOK, map[string]string{"key": "value"})
	}).Methods("GET")

	router.HandleFunc("/api/v1/tasks/results/{task_uuid}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		taskUUID := vars["task_uuid"]
		t, err := database.TaskResultsByUUID(dbConn, taskUUID)
		if err != nil {
			log.Println(err)
		}
		ren.JSON(w, http.StatusOK, t)
	}).Methods("GET")

	router.HandleFunc("/api/v1/tasks/{task_uuid}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		taskUUID := vars["task_uuid"]
		t, err := database.TaskByUUID(dbConn, taskUUID)
		if err != nil {
			log.Println(err)
		}
		ren.JSON(w, http.StatusOK, t)
	}).Methods("GET")

	router.HandleFunc("/api/v1/tasks/running", func(w http.ResponseWriter, r *http.Request) {
		t, err := database.TasksRunning(dbConn)
		if err != nil {
			log.Println(err)
		}
		ren.JSON(w, http.StatusOK, t)
	}).Methods("GET")

	router.HandleFunc("/api/v1/schedules", func(w http.ResponseWriter, r *http.Request) {
		results, err := ShowSchedules(conf)
		if err != nil {
			log.Println(err)
		}
		fmt.Println(results)
		ren.JSON(w, http.StatusOK, map[string]string{"sent": "value"})
	}).Methods("GET")

	router.HandleFunc("/api/v1/schedules/{task_id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		taskID := vars["task_id"]
		taskChan <- taskID
		ren.JSON(w, http.StatusOK, map[string]string{"sent": taskID})
	}).Methods("GET")

	router.HandleFunc("/api/v1/schedules", func(w http.ResponseWriter, r *http.Request) {
		ren.JSON(w, http.StatusOK, map[string]string{"sent": "value"})
	}).Methods("POST")

	router.HandleFunc("/api/v1/schedules/delete/{schedule_id}", func(w http.ResponseWriter, r *http.Request) {
		ren.JSON(w, http.StatusOK, map[string]string{"sent": "value"})
	}).Methods("DELETE")

	//
	// Page routes
	//
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ren.HTML(w, http.StatusOK, "index", nil)
	}).Methods("GET")
	router.HandleFunc("/reports", func(w http.ResponseWriter, r *http.Request) {
		ren.HTML(w, http.StatusOK, "reports", nil)
	}).Methods("GET")
	router.HandleFunc("/analytics", func(w http.ResponseWriter, r *http.Request) {
		ren.HTML(w, http.StatusOK, "analytics", nil)
	}).Methods("GET")
	router.HandleFunc("/export", func(w http.ResponseWriter, r *http.Request) {
		ren.HTML(w, http.StatusOK, "export", nil)
	}).Methods("GET")
	router.HandleFunc("/races", func(w http.ResponseWriter, r *http.Request) {
		ren.HTML(w, http.StatusOK, "races", nil)
	}).Methods("GET")

	n.UseHandler(router)
	n.Run(conf.Raceway.AppPort)
}
