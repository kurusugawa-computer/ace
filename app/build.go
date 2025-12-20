package app

import (
	"errors"
	"strings"
	"text/template"

	"github.com/kurusugawa-computer/ace/agents"
)

const DefaultApprovalPolicy = "never"
const DefaultSandbox = "read-only"
const DefaultTimeoutSec = 1800

func (app *App) buildAgent(agentName string) (*agents.Agent, error) {
	// エージェントのConfigを取得
	agentConfig, ok := app.config.Agents[agentName]
	if !ok {
		return nil, errors.New("no such agent: " + agentName)
	}

	// 共通CodexConfigをベースにエージェントのCodexConfigをマージ
	codexConfig := app.config.Config.Clone()
	for key, value := range agentConfig.Config {
		codexConfig[key] = value
	}

	// ApprovalPolicy の解決
	approvalPolicy := agentConfig.ApprovalPolicy
	if approvalPolicy == "" {
		approvalPolicy = DefaultApprovalPolicy
	}

	// Sandbox の解決
	sandbox := agentConfig.Sandbox
	if sandbox == "" {
		sandbox = DefaultSandbox
	}

	// base instructions の送信可否
	useBaseInstructions := true
	if value, ok := codexConfig["use_base_instructions"]; ok {
		enabled, ok := value.(bool)
		if !ok {
			return nil, errors.New("use_base_instructions must be a boolean")
		}
		useBaseInstructions = enabled
		delete(codexConfig, "use_base_instructions")
	}

	// MCP Servers の解決
	for mcpServerName, mcpServerConfig := range agentConfig.MCPServers {
		codexConfig["mcp_servers."+mcpServerName] = mcpServerConfig
	}

	// サブエージェントの解決
	subAgents := make([]*agents.SubAgentConfig, 0, len(agentConfig.SubAgents))
	for _, subAgentName := range agentConfig.SubAgents {
		subAgentConfig, ok := app.config.Agents[subAgentName]
		if !ok {
			return nil, errors.New("no such agent: " + subAgentName)
		}
		timeoutSec := subAgentConfig.TimeoutSec
		if timeoutSec == 0 {
			timeoutSec = DefaultTimeoutSec
		}
		subAgents = append(subAgents, &agents.SubAgentConfig{
			Name:       subAgentConfig.Name,
			TimeoutSec: timeoutSec,
		})
	}

	// vars の値を適用
	description := agentConfig.Description
	instruction := agentConfig.Instruction
	if app.config.Vars != nil {
		descriptionTemplate, err := template.New("description").Parse(agentConfig.Description)
		if err != nil {
			return nil, err
		}
		descriptionBuilder := &strings.Builder{}
		if err := descriptionTemplate.Execute(descriptionBuilder, app.config.Vars); err != nil {
			return nil, err
		}

		instructionTemplate, err := template.New("instruction").Parse(agentConfig.Instruction)
		if err != nil {
			return nil, err
		}
		instructionBuilder := &strings.Builder{}
		if err := instructionTemplate.Execute(instructionBuilder, app.config.Vars); err != nil {
			return nil, err
		}

		description = descriptionBuilder.String()
		instruction = instructionBuilder.String()
	}

	// エージェントのビルド
	agent, err := agents.Build(&agents.Config{
		Name:           agentConfig.Name,
		Description:    description,
		Instruction:    instruction,
		PromptTemplate: agentConfig.PromptTemplate,
		InputSchema:    agentConfig.InputSchema,
		OutputSchema:   agentConfig.OutputSchema,
		ApprovalPolicy: approvalPolicy,
		Sandbox:        sandbox,
		Config:              codexConfig,
		UseBaseInstructions: useBaseInstructions,
		SubAgents:           subAgents,
	})
	if err != nil {
		return nil, err
	}

	return agent, nil
}
