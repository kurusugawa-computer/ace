package app

import (
	"context"
	"errors"

	ag "github.com/kurusugawa-computer/ace/app/agent"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func (app *App) RunMCPServer(agentName string, workdir string) error {
	// エージェントのConfigを取得
	agentConfig, ok := app.config.Agents[agentName]
	if !ok {
		return errors.New("no such agent: " + agentName)
	}

	// エージェントのビルド
	agent, err := ag.Build(agentConfig)
	if err != nil {
		return err
	}

	// MCP Serverを構築
	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    agent.Name,
			Version: "v1.0.0",
		},
		&mcp.ServerOptions{},
	)
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:         agent.Name,
			Description:  agent.Description,
			InputSchema:  agent.InputSchema,
			OutputSchema: agent.OutputSchema,
		},
		func(ctx context.Context, request *mcp.CallToolRequest, input map[string]any) (*mcp.CallToolResult, any, error) {
			output, err := agent.Run(
				workdir,
				input,
				&ag.RunConfig{
					APIKey:                  app.apiKey,
					SubAgentMCPServerConfig: app.subAgentMCPServerConfig,
					LogLevel:                app.logLevel,
					LogWriter:               app.logWriter,
				},
			)
			if err != nil {
				return nil, nil, err
			}
			return nil, output, nil
		},
	)

	// MCP Serverを起動
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		return err
	}

	return nil
}
