package controllers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"

	"github.com/briandowns/raceway/config"
	"github.com/briandowns/raceway/database"
)

// TasksHandler
func TasksHandler(ren *render.Render, conf *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db, err := database.NewDatabase(conf)
		if err != nil {
			log.Fatalln(err)
			return
		}
		ren.JSON(w, http.StatusOK, db.GetTasks())
	}
}

// TasksByUUIDHandler
func TasksByUUIDHandler(ren *render.Render, conf *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		uuid := vars["task_uuid"]
		db, err := database.NewDatabase(conf)
		if err != nil {
			log.Fatalln(err)
			return
		}
		ren.JSON(w, http.StatusOK, db.TaskByUUID(uuid))
	}
}

// TasksResultsByUUIDHandler
func TasksResultsByUUIDHandler(ren *render.Render, conf *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		uuid := vars["task_uuid"]
		db, err := database.NewDatabase(conf)
		if err != nil {
			log.Fatalln(err)
			return
		}
		ren.JSON(w, http.StatusOK, db.TaskResultsByUUID(uuid))
	}
}

// TasksRunningHandler
func TasksRunningHandler(ren *render.Render, conf *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db, err := database.NewDatabase(conf)
		if err != nil {
			log.Fatalln(err)
			return
		}
		ren.JSON(w, http.StatusOK, db.TasksRunning())
	}
}
