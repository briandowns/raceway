package controllers

var Frontend = "/"

const (
	// APIBase is the base path for API access
	APIBase = "/api/v1/"

	DeploymentsPath       = APIBase + "deployments"
	DeploymentsByNamePath = DeploymentsPath + "/{name}"

	ScenariosPath       = APIBase + "scenarios"
	ScenariosByNamePath = ScenariosPath + "/{name}"

	TasksPath        = APIBase + "tasks"
	TasksByUUIDPath  = TasksPath + "/{uuid}"
	TasksResultPath  = TasksPath + "/results/{uuid}"
	TasksRunningPath = TasksPath + "/running"

	SchedulesPath           = APIBase + "schedules"
	SchedulesDeletePath     = SchedulesPath + "/{uuid}"
	SchedulesByTaskUUIDPath = SchedulesPath + "{task_uuid}"
)
