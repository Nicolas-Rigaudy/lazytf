package terraform

import (
	"encoding/json"
	"os/exec"
)

type PlanSummary struct {
	ResourceChanges []ResourceChange `json:"resource_changes"`
}

type ResourceChange struct {
	Address string `json:"address"`
	Change  Change `json:"change"`
}

type Change struct {
	Actions []string               `json:"actions"`
	Before  map[string]interface{} `json:"before"`
	After   map[string]interface{} `json:"after"`
}

func ParsePlanJSON(jsonData []byte) (*PlanSummary, error) {
	var plan PlanSummary
	err := json.Unmarshal(jsonData, &plan)
	if err != nil {
		return nil, err
	}

	return &plan, nil
}

func RunShowPlanJSON(projectPath string, planFilePath string) ([]byte, error) {
	cmd := exec.Command("terraform", "show", "-json", planFilePath)
	cmd.Dir = projectPath
	return cmd.Output()
}
