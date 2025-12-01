package agent

import "github.com/google/jsonschema-go/jsonschema"

type Config struct {
	Name           string                        `yaml:"-"`
	Description    string                        `yaml:"description"`
	Instruction    string                        `yaml:"instruction"`
	PromptTemplate string                        `yaml:"prompt_template"`
	InputSchema    map[string]*jsonschema.Schema `yaml:"input_schema"`
	OutputSchema   map[string]*jsonschema.Schema `yaml:"output_schema"`
	ApprovalPolicy string                        `yaml:"approval_policy"` // untrusted, on-failure, never
	Sandbox        string                        `yaml:"sandbox"`         // read-only, workspace-write, danger-full-access
	Config         map[string]any                `yaml:"config"`
	SubAgents      []string                      `yaml:"sub_agents"`
}
