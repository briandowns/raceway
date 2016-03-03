package database

// Deployments is the model for the rally.deployments table
type Deployments struct {
	ID                    int    `sql:"auto_increment" gorm:"column:id" gorm:"primary_key" json:"id"`
	CreatedAt             string `gorm:"column:created_at" json:"created_at"`
	UpdatedAt             string `gorm:"column:updated_at" json:"updated_at"`
	UUID                  string `gorm:"column:uuid" json:"uuid"`
	ParentUUID            []byte `gorm:"column:parent_uuid" json:"parent_uuid"`
	Name                  string `gorm:"column:name" json:"name"`
	StartedAt             string `gorm:"column:started_at" json:"started_at"`
	CompletedAt           string `gorm:"column:completed_at" json:"completed_at"`
	EnumDeploymentsStatus string `gorm:"column:enum_deployments_status" json:"enum_deployments_status"`
}

// Resources is the model for the rally.Resources table
type Resources struct {
	ID             int    `sql:"auto_increment" gorm:"column:id" gorm:"primary_key" json:"id"`
	CreatedAt      string `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      string `gorm:"column:updated_at" json:"updated_at"`
	ProviderName   string `gorm:"column:provider_name" json:"provider_name"`
	Type           string `gorm:"column:type" json:"type"`
	Info           string `gorm:"column:info" json:"info"`
	DeploymentUUID string `gorm:"column:deployment_uuid" json:"deployment_uuid"`
}

// Tasks is the model for the rally.tasks table
type Tasks struct {
	ID              int    `sql:"auto_increment" gorm:"column:id" gorm:"primary_key" json:"id"`
	CreatedAt       string `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       string `gorm:"column:updated_at" json:"updated_at"`
	UUID            string `gorm:"column:uuid" json:"uuid"`
	Status          string `gorm:"column:status" json:"status"`
	VerificationLog string `db:"verification_log"`
	Tag             string `gorm:"column:tag" json:"tag"`
	DeploymentUUID  string `gorm:"column:deployment_uuid" json:"deployment_uuid"`
}

// TaskResults is the model for the rally.task_results table
type TaskResults struct {
	ID        int    `sql:"auto_increment" gorm:"column:id" gorm:"primary_key" json:"id"`
	CreatedAt string `gorm:"column:created_at" json:"created_at"`
	UpdatedAt string `gorm:"column:updated_at" json:"updated_at"`
	Key       string `gorm:"column:key" json:"key"`
	Data      string `gorm:"column:data" json:"data"`
	TaskUUID  string `gorm:"column:task_uuid" json:"task_uuid"`
}

// VerificationResults is the model for the rally.verification_results table
type VerificationResults struct {
	ID               int    `sql:"auto_increment" gorm:"column:id" gorm:"primary_key" json:"id"`
	CreatedAt        string `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        string `gorm:"column:updated_at" json:"updated_at"`
	VerificationUUID string `gorm:"column:verification_uuid" json:"verification_uuid"`
	Data             string `gorm:"column:data" json:"data"`
}

// Verifications is the model for the rally.verifications table
type Verifications struct {
	ID             int     `sql:"auto_increment" gorm:"column:id" gorm:"primary_key" json:"id"`
	CreatedAt      string  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      string  `gorm:"column:updated_at" json:"updated_at"`
	UUID           string  `gorm:"column:uuid" json:"uuid"`
	DeploymentUUID string  `gorm:"column:deployment_uuid" json:"deployment_uuid"`
	Status         string  `gorm:"column:status" json:"status"`
	SetName        string  `gorm:"column:set_name" json:"set_name"`
	Tests          int     `gorm:"column:tests" json:"tests"`
	Errors         int     `gorm:"column:errors" json:"errors"`
	Failures       int     `gorm:"column:failures" json:"failures"`
	Time           float64 `gorm:"column:time" json:"time"`
}

// Workers is the model for the rally.workers table
type Workers struct {
	ID        int    `sql:"auto_increment" gorm:"column:id" gorm:"primary_key" json:"id"`
	CreatedAt string `gorm:"column:created_at" json:"created_at"`
	UpdatedAt string `gorm:"column:updated_at" json:"updated_at"`
	Hostname  string `gorm:"column:hostname" json:"hostname"`
}
