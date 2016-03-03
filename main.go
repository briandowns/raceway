package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/pborman/uuid"

	"github.com/codegangsta/negroni"
	"github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"

	"github.com/briandowns/raceway/config"
	"github.com/briandowns/raceway/database"
)

var taskChan = make(chan string, 0)
var Conf *config.Config

var signalsChan = make(chan os.Signal, 1)

// Task holds task details once initialized
type Task struct{}

func main() {
	signal.Notify(signalsChan, os.Interrupt)

	go func() {
		for sig := range signalsChan {
			log.Printf("Exiting... %v\n", sig)
			signalsChan = nil
			os.Exit(1)
		}
	}()

	// Source in the configuration
	conf, err := config.Load("./config.json")
	if err != nil {
		log.Fatalln(err)
	}

	db, err := database.NewDatabase(conf)
	if err != nil {
		log.Fatalln(err)
	}

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
		ren.JSON(w, http.StatusOK, db.GetDeployments())
	}).Methods("GET")

	router.HandleFunc("/api/v1/deployments/{name}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		ren.JSON(w, http.StatusOK, db.DeploymentsByName(name))
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
		ren.JSON(w, http.StatusOK, db.GetTasks())
	}).Methods("GET")

	router.HandleFunc("/api/v1/tasks/start", func(w http.ResponseWriter, r *http.Request) {
		ren.JSON(w, http.StatusOK, map[string]string{"key": "value"})
	}).Methods("GET")

	router.HandleFunc("/api/v1/tasks/results/{task_uuid}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		taskUUID := vars["task_uuid"]
		ren.JSON(w, http.StatusOK, db.TaskResultsByUUID(taskUUID))
	}).Methods("GET")

	router.HandleFunc("/api/v1/tasks/{task_uuid}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		taskUUID := vars["task_uuid"]
		ren.JSON(w, http.StatusOK, db.TaskByUUID(taskUUID))
	}).Methods("GET")

	router.HandleFunc("/api/v1/tasks/running", func(w http.ResponseWriter, r *http.Request) {
		ren.JSON(w, http.StatusOK, db.TasksRunning())
	}).Methods("GET")

	router.HandleFunc("/api/v1/schedules", func(w http.ResponseWriter, r *http.Request) {
		results, err := ShowSchedules()
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
