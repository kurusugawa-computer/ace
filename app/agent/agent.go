package agent

import (
	"html/template"

	"github.com/google/jsonschema-go/jsonschema"
)

type Agent struct {
	Name           string
	Description    string
	Instruction    string
	PromptTemplate *template.Template
	InputSchema    *jsonschema.Schema
	OutputSchema   *jsonschema.Schema
	ApprovalPolicy string
	Sandbox        string
	Config         map[string]any
	SubAgents      []string
}
