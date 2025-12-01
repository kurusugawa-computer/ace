package cli

import (
	"os"
	"path/filepath"

	"github.com/kurusugawa-computer/ace/agents"
	"github.com/urfave/cli/v3"
)

func New(appName string, version string, title string) *cli.Command {
	return &cli.Command{
		Name:                  appName,
		Usage:                 title,
		ArgsUsage:             " ",
		Version:               version,
		Description:           `Define, run, and manage AI agents using YAML configuration and JSON Schema.`,
		EnableShellCompletion: true,
		HideHelpCommand:       true,
		HideVersion:           false,
		Commands: []*cli.Command{
			exec(appName, version),
			mcp(appName, version),
			setup(appName, version),
		},
	}
}

type subCommand func(appName string, version string) *cli.Command

// サブエージェントを実行するMCP Serverの起動方法を返す関数を返す関数
func subAgentMCPServerConfig(configPath string, workdir string, apiKey string) func(subAgent *agents.SubAgent) (map[string]any, error) {
	return func(subAgent *agents.SubAgent) (map[string]any, error) {
		// 設定ファイルの絶対パスを取得
		configAbsPath, err := filepath.Abs(configPath)
		if err != nil {
			return nil, err
		}

		// 作業ディレクトリの絶対パスを取得
		workdirAbsPath, err := filepath.Abs(workdir)
		if err != nil {
			return nil, err
		}

		// サブエージェント MCP Server 用の Config を構築
		config := map[string]any{
			"command": os.Args[0],
			"args": []string{
				"mcp-server",
				"--config",
				configAbsPath,
				"--workdir",
				workdirAbsPath,
				subAgent.Name,
			},
			"env": map[string]any{
				"OPENAI_API_KEY": apiKey,
			},
			"startup_timeout_sec": 30,
			"tool_timeout_sec":    subAgent.TimeoutSec,
		}

		return config, nil
	}
}
