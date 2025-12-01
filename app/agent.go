package app

import (
	"github.com/kurusugawa-computer/ace/agents"
)

func (app *App) RunAgent(agentName string, workdir string, input map[string]any) (any, error) {
	// エージェントのビルド
	agent, err := app.buildAgent(agentName)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return output, nil
}
