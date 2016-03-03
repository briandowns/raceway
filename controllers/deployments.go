package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/unrolled/render"

	"github.com/briandowns/raceway/config"
	"github.com/briandowns/raceway/database"
)

// NewDeploymentParams holds the decoded JSON from a POST to create
// a new deployment
type NewDeploymentParams struct {
	FromENV bool   `json:"from_env"`
	Name    string `json:"name"`
}

// DeploymentsHandler
func DeploymentsHandler(ren *render.Render, conf *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db, err := database.NewDatabase(conf)
		if err != nil {
			log.Fatalln(err)
			return
		}
		ren.JSON(w, http.StatusOK, db.GetDeployments())
	}
}

// DeploymentsByNameHandler
func DeploymentsByNameHandler(ren *render.Render, conf *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db, err := database.NewDatabase(conf)
		if err != nil {
			log.Fatalln(err)
			return
		}
		var dp NewDeploymentParams
		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(&dp)
		if err != nil {
			log.Println(err)
			return
		}
		ren.JSON(w, http.StatusOK, db.DeploymentsByName(dp.Name))
	}
}

// NewDeploymentsHandler
func NewDeploymentsHandler(ren *render.Render, conf *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ren.JSON(w, http.StatusOK, map[string]interface{}{"deployment": "new"})
	}
}
