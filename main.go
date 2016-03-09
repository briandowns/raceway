package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/codegangsta/negroni"
	"github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	"github.com/gorilla/mux"
	"github.com/pborman/uuid"
	"github.com/thoas/stats"
	"github.com/unrolled/render"

	"github.com/briandowns/raceway/config"
	"github.com/briandowns/raceway/controllers"
	"github.com/briandowns/raceway/scheduler"
)

var taskChan = make(chan string, 0)
var Conf *config.Config

var signalsChan = make(chan os.Signal, 1)

var configFileFlag string

func init() {
	flag.StringVar(&configFileFlag, "config", "", "path to raceway config file")
}

func main() {
	flag.Parse()

	signal.Notify(signalsChan, os.Interrupt)

	go func() {
		for sig := range signalsChan {
			log.Printf("Exiting... %v\n", sig)
			signalsChan = nil
			os.Exit(1)
		}
	}()

	// don't go any further if no config is supplied
	if configFileFlag == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Source in the configuration
	conf, err := config.Load(configFileFlag)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("%+v\n", conf)
	// Start the scheduler
	scheduler.StartScheduler(conf)

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
	router.HandleFunc(controllers.DeploymentsPath, controllers.DeploymentsHandler(ren, conf)).Methods("GET")

	router.HandleFunc(controllers.DeploymentsByNamePath, controllers.DeploymentsByNameHandler(ren, conf)).Methods("GET")

	router.HandleFunc(controllers.DeploymentsPath, controllers.NewDeploymentsHandler(ren, conf)).Methods("POST")

	router.HandleFunc(controllers.ScenariosPath, controllers.ScenariosHandler(ren, conf)).Methods("GET")

	router.HandleFunc(controllers.ScenariosByNamePath, controllers.ScenariosByNameHandler(ren, conf)).Methods("GET")

	router.HandleFunc(controllers.TasksPath, controllers.TasksHandler(ren, conf)).Methods("GET")

	router.HandleFunc(controllers.TasksResultPath, controllers.TasksResultsByUUIDHandler(ren, conf)).Methods("GET")

	router.HandleFunc(controllers.TasksByUUIDPath, controllers.TasksByUUIDHandler(ren, conf)).Methods("GET")

	router.HandleFunc(controllers.TasksRunningPath, controllers.TasksRunningHandler(ren, conf)).Methods("GET")

	router.HandleFunc(controllers.SchedulesPath, controllers.SchedulesHandler(ren, conf)).Methods("GET")

	router.HandleFunc(controllers.SchedulesByTaskUUIDPath, controllers.SchedulesByTaskIDHandler(ren, conf)).Methods("GET")

	router.HandleFunc(controllers.SchedulesPath, controllers.NewScheduleHandler(ren, conf)).Methods("POST")

	router.HandleFunc(controllers.SchedulesDeletePath, controllers.DeleteScheduleHandler(ren, conf)).Methods("DELETE")

	n.Use(statsMiddleware)
	n.UseHandler(router)
	n.Run(conf.Raceway.AppPort)
}
