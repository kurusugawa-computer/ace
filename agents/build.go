package agents

import (
	"html/template"

	"github.com/google/jsonschema-go/jsonschema"
)

func Build(config *Config) (*Agent, error) {
	// 入出力のスキーマ定義
	inputSchema := &jsonschema.Schema{
		Type:                 "object",
		Required:             []string{},
		Properties:           map[string]*jsonschema.Schema{},
		AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
	}
	for name, schema := range config.InputSchema {
		inputSchema.Required = append(inputSchema.Required, name)
		inputSchema.Properties[name] = schema
	}

	outputSchema := &jsonschema.Schema{
		Type:                 "object",
		Required:             []string{},
		Properties:           map[string]*jsonschema.Schema{},
		AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
	}
	for name, schema := range config.OutputSchema {
		outputSchema.Required = append(outputSchema.Required, name)
		outputSchema.Properties[name] = schema
	}

	// プロンプトテンプレートのビルド
	promptTemplate, err := template.New("prompt").Parse(config.PromptTemplate)
	if err != nil {
		return nil, err
	}

	// サブエージェント
	subAgents := make([]*SubAgent, 0, len(config.SubAgents))
	for _, subAgentConfig := range config.SubAgents {
		subAgents = append(subAgents, &SubAgent{
			Name:       subAgentConfig.Name,
			TimeoutSec: subAgentConfig.TimeoutSec,
		})
	}

	// 構築したエージェントを返す
	agent := &Agent{
		Name:           config.Name,
		Description:    config.Description,
		Instruction:    config.Instruction,
		PromptTemplate: promptTemplate,
		InputSchema:    inputSchema,
		OutputSchema:   outputSchema,
		ApprovalPolicy: config.ApprovalPolicy,
		Sandbox:        config.Sandbox,
		Config:         config.Config,
		SubAgents:      subAgents,
	}
	return agent, nil
}
