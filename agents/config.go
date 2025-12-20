package agents

import (
	"encoding/json"

	"github.com/google/jsonschema-go/jsonschema"
)

type Config struct {
	Name           string
	Description    string
	Instruction    string
	PromptTemplate string
	InputSchema    map[string]*jsonschema.Schema
	OutputSchema   map[string]*jsonschema.Schema
	ApprovalPolicy string // untrusted, on-failure, never
	Sandbox        string // read-only, workspace-write, danger-full-access
	Config         CodexConfig
	UseBaseInstructions bool
	SubAgents      []*SubAgentConfig
}

type SubAgentConfig struct {
	Name       string
	TimeoutSec int
}

type CodexConfig map[string]any

func (codexConfig CodexConfig) Clone() CodexConfig {
	copied := map[string]any{}
	if codexConfig != nil {
		jsonConfig, _ := json.Marshal(codexConfig)
		_ = json.Unmarshal(jsonConfig, &copied)
	}
	return copied
}
