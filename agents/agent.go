package agents

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
	Config         CodexConfig
	SubAgents      []*SubAgent
}

type SubAgent struct {
	Name       string
	TimeoutSec int
}
