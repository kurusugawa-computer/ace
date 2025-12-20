package app

import (
	"io"

	"github.com/kurusugawa-computer/ace/agents"
)

type App struct {
	config                  *Config
	codexExecutablePath     string // Codex の実行パス
	apiKey                  string
	subAgentMCPServerConfig func(subAgent *agents.SubAgent) (map[string]any, error)

	logWriter io.Writer
	logLevel  string // error, warn, info, debug, trace, off
}

type AppOption func(*App)

func New(config *Config, codexExecutablePath string, apiKey string, subAgentMCPServerConfig func(subAgent *agents.SubAgent) (map[string]any, error), options ...AppOption) *App {
	app := &App{
		config:                  config,
		codexExecutablePath:     codexExecutablePath,
		apiKey:                  apiKey, // codex login でログイン済みなら空文字列
		subAgentMCPServerConfig: subAgentMCPServerConfig,
	}

	for _, option := range options {
		option(app)
	}

	return app
}

func WithLogger(logWriter io.Writer, logLevel string) AppOption {
	return func(app *App) {
		app.logWriter = logWriter
		app.logLevel = logLevel
	}
}
