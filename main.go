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
	"github.com/thoas/stats"
	"github.com/unrolled/render"

	"github.com/briandowns/raceway/config"
	"github.com/briandowns/raceway/controllers"
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

	statsMiddleware := stats.New()

	router := mux.NewRouter()

	// Frontend
	router.HandleFunc(controllers.Frontend, controllers.FrontendHandler()).Methods("GET")

	// API routes
	router.HandleFunc(controllers.DeploymentsPath, func(w http.ResponseWriter, r *http.Request) {
		ren.JSON(w, http.StatusOK, db.GetDeployments())
	}).Methods("GET")

	router.HandleFunc(controllers.DeploymentsByNamePath, func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		ren.JSON(w, http.StatusOK, db.DeploymentsByName(name))
	}).Methods("GET")

	router.HandleFunc(controllers.DeploymentsPath, func(w http.ResponseWriter, r *http.Request) {
		dp := &NewDeploymentParams{}
		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(&dp)
		if err != nil {
			log.Println(err)
		}
		ren.JSON(w, http.StatusOK, dp)
	}).Methods("POST")

	router.HandleFunc(controllers.ScenariosPath, func(w http.ResponseWriter, r *http.Request) {
		ren.JSON(w, http.StatusOK, Scenarios(conf.Scenarios.ScenarioDir, conf.Scenarios.ScenarioFormat))
	}).Methods("GET")

	router.HandleFunc(controllers.ScenariosByNamePath, func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		scenarioName := vars["name"]
		ren.JSON(w, http.StatusOK, ScenarioContent(conf.Scenarios.ScenarioDir, scenarioName))
	}).Methods("GET")

	router.HandleFunc(controllers.TasksPath, func(w http.ResponseWriter, r *http.Request) {
		ren.JSON(w, http.StatusOK, db.GetTasks())
	}).Methods("GET")

	router.HandleFunc(controllers.TasksResultPath, func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		taskUUID := vars["task_uuid"]
		ren.JSON(w, http.StatusOK, db.TaskResultsByUUID(taskUUID))
	}).Methods("GET")

	router.HandleFunc(controllers.TasksByUUIDPath, func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		taskUUID := vars["task_uuid"]
		ren.JSON(w, http.StatusOK, db.TaskByUUID(taskUUID))
	}).Methods("GET")

	router.HandleFunc(controllers.TasksRunningPath, func(w http.ResponseWriter, r *http.Request) {
		ren.JSON(w, http.StatusOK, db.TasksRunning())
	}).Methods("GET")

	router.HandleFunc(controllers.SchedulesPath, func(w http.ResponseWriter, r *http.Request) {
		results, err := ShowSchedules()
		if err != nil {
			log.Println(err)
		}
		fmt.Println(results)
		ren.JSON(w, http.StatusOK, map[string]string{"sent": "value"})
	}).Methods("GET")

	router.HandleFunc(controllers.SchedulesByTaskUUIDPath, func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		taskID := vars["task_id"]
		taskChan <- taskID
		ren.JSON(w, http.StatusOK, map[string]string{"sent": taskID})
	}).Methods("GET")

	router.HandleFunc(controllers.SchedulesPath, func(w http.ResponseWriter, r *http.Request) {
		ren.JSON(w, http.StatusOK, map[string]string{"sent": "value"})
	}).Methods("POST")

	router.HandleFunc(controllers.SchedulesDeletePath, func(w http.ResponseWriter, r *http.Request) {
		ren.JSON(w, http.StatusOK, map[string]string{"sent": "value"})
	}).Methods("DELETE")

	n.Use(statsMiddleware)
	n.UseHandler(router)
	n.Run(conf.Raceway.AppPort)
}
