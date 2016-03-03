package controllers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"

	"github.com/briandowns/raceway/config"
	"github.com/briandowns/raceway/scheduler"
)

// SchedulesHandler
func SchedulesHandler(ren *render.Render, conf *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		results, err := scheduler.ShowSchedules(conf)
		if err != nil {
			log.Println(err)
			return
		}
		ren.JSON(w, http.StatusOK, map[string]interface{}{"schedules": results})
	}
}

// SchedulesByTaskIDHandler
func SchedulesByTaskIDHandler(ren *render.Render, conf *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		taskID := vars["task_id"]
		ren.JSON(w, http.StatusOK, map[string]string{"sent": taskID})
	}
}

// NewScheduleHandler
func NewScheduleHandler(ren *render.Render, conf *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ren.JSON(w, http.StatusOK, map[string]string{"schedule": "anything for right now..."})
	}
}

// DeleteScheduleHandler
func DeleteScheduleHandler(ren *render.Render, conf *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ren.JSON(w, http.StatusOK, map[string]string{"schedule": "anything for right now..."})
	}
}
