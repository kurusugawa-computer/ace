package app

import (
	"errors"

	ag "github.com/kurusugawa-computer/ace/app/agent"
)

func (app *App) RunAgent(agentName string, workdir string, input map[string]any) (any, error) {
	// エージェントのConfigを取得
	agentConfig, ok := app.config.Agents[agentName]
	if !ok {
		return nil, errors.New("no such agent: " + agentName)
	}

	// エージェントのビルド
	agent, err := ag.Build(agentConfig)
	if err != nil {
		return nil, err
	}

	// エージェントの実行
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
		return nil, err
	}

	return output, nil
}
