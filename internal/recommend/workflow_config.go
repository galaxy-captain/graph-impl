package recommend

type WorkflowMap struct {
	Workflows map[string]WorkflowConfig `json:"workflows"`
}

type WorkflowConfig struct {
	Sections map[string]WorkflowSectionConfig `json:"sections"`
}

type WorkflowSectionConfig struct {
	Capacity int                           `json:"capacity"`
	Nodes    map[string]WorkflowNodeConfig `json:"nodes"`
}

type WorkflowNodeConfig struct {
	ClassName  string   `json:"class_name"`
	InputNodes []string `json:"input_nodes"`
	Count      int      `json:"count"`
}
