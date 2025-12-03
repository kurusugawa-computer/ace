package app

import (
	"context"

	"github.com/kurusugawa-computer/ace/agents"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func (app *App) RunMCPServer(agentName string, workdir string) error {
	// エージェントのビルド
	agent, err := app.buildAgent(agentName)
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
			// vars の値を展開する
			if app.config.Vars != nil {
				for key, value := range app.config.Vars {
					if _, ok := input[key]; ok {
						continue
					}
					input[key] = value
				}
			}

			// エージェントの実行
			output, err := agent.Run(
				workdir,
				input,
				&agents.RunConfig{
					APIKey:                  app.apiKey,
					SubagentMCPServerConfig: app.subAgentMCPServerConfig,
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
