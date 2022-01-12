package client

// ResProcessDefinition a JSON object corresponding to the ProcessDefinition interface in the engine
type ResProcessDefinition struct {
	// The id of the process definition
	Id string `json:"id"`
	// The key of the process definition, i.e., the id of the BPMN 2.0 XML process definition
	Key string `json:"key"`
	// The category of the process definition
	Category string `json:"category"`
	// The description of the process definition
	Description string `json:"description"`
	// The name of the process definition
	Name string `json:"name"`
	// The version of the process definition that the engine assigned to it
	Version int `json:"Version"`
	// The file name of the process definition
	Resource string `json:"resource"`
	// The deployment id of the process definition
	DeploymentId string `json:"deploymentId"`
	// The file name of the process definition diagram, if it exists
	Diagram string `json:"diagram"`
	// A flag indicating whether the definition is suspended or not
	Suspended bool `json:"suspended"`
	// The tenant id of the process definition
	TenantId string `json:"tenantId"`
	// The version tag of the process definition
	VersionTag string `json:"versionTag"`
	// History time to live value of the process definition. Is used within History cleanup
	HistoryTimeToLive int `json:"historyTimeToLive"`
	// A flag indicating whether the process definition is startable in Tasklist or not
	StartableInTasklist bool `json:"startableInTasklist"`
}
