package terraform

import (
	"encoding/json"
	"os"
	"os/exec"
)

type PlanSummary struct {
	ResourceChanges []ResourceChange `json:"resource_changes"`
}

type ResourceChange struct {
	Address       string `json:"address"`
	ModuleAddress string `json:"module_address"` // empty for root module resources
	Change        Change `json:"change"`
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

func RunShowPlanJSON(tool Tool, projectPath string, planFilePath string) ([]byte, error) {
	cmd := exec.Command(Binary(tool.Bin), "show", "-json", planFilePath)
	if len(tool.Env) > 0 {
		cmd.Env = append(os.Environ(), tool.Env...)
	}
	cmd.Dir = projectPath
	return cmd.Output()
}
