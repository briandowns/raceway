package database

// Deployments is the model for the rally.deployments table
type Deployments struct {
	CreatedAt             string `db:"created_at"`
	UpdatedAt             string `db:"updated_at"`
	ID                    int    `db:"id"`
	UUID                  string `db:"uuid"`
	ParentUUID            []byte `db:"parent_uuid"`
	Name                  string `db:"name"`
	StartedAt             string `db:"started_at"`
	CompletedAt           string `db:"completed_at"`
	EnumDeploymentsStatus string `db:"enum_deployments_status"`
}

// Resources is the model for the rally.Resources table
type Resources struct {
	CreatedAt      string `db:"created_at"`
	UpdatedAt      string `db:"updated_at"`
	ID             int    `db:"id"`
	ProviderName   string `db:"provider_name"`
	Type           string `db:"type"`
	Info           string `db:"info"`
	DeploymentUUID string `db:"deployment_uuid"`
}

// Tasks is the model for the rally.tasks table
type Tasks struct {
	CreatedAt       string `db:"created_at"`
	UpdatedAt       string `db:"updated_at"`
	ID              int    `db:"id"`
	UUID            string `db:"uuid"`
	Status          string `db:"status"`
	VerificationLog string `db:"verification_log"`
	Tag             string `db:"tag"`
	DeploymentUUID  string `db:"deployment_uuid"`
}

// TaskResults is the model for the rally.task_results table
type TaskResults struct {
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
	ID        int    `db:"id"`
	Key       string `db:"key"`
	Data      string `db:"data"`
	TaskUUID  string `db:"task_uuid"`
}

// VerificationResults is the model for the rally.verification_results table
type VerificationResults struct {
	CreatedAt        string `db:"created_at"`
	UpdatedAt        string `db:"updated_at"`
	ID               int    `db:"id"`
	VerificationUUID string `db:"verification_uuid"`
	Data             string `db:"data"`
}

// Verifications is the model for the rally.verifications table
type Verifications struct {
	CreatedAt      string  `db:"created_at"`
	UpdatedAt      string  `db:"updated_at"`
	ID             int     `db:"id"`
	UUID           string  `db:"uuid"`
	DeploymentUUID string  `db:"deployment_uuid"`
	Status         string  `db:"status"`
	SetName        string  `db:"set_name"`
	Tests          int     `db:"tests"`
	Errors         int     `db:"errors"`
	Failures       int     `db:"failures"`
	Time           float64 `db:"time"`
}

// Workers is the model for the rally.workers table
type Workers struct {
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
	ID        int    `db:"id"`
	Hostname  string `db:"hostname"`
}
