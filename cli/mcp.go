package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/kurusugawa-computer/ace/app"
	"github.com/urfave/cli/v3"
)

var _ subCommand = mcp

func mcp(appName string, version string) *cli.Command {
	return &cli.Command{
		Name:      "mcp-server",
		Aliases:   []string{},
		Usage:     "Start Ace as an MCP server to serve multiple agents.",
		ArgsUsage: "AGENT_NAME",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "",
				Value:   "agent.yaml",
			},
			&cli.StringFlag{
				Name:    "workdir",
				Aliases: []string{"w"},
				Usage:   "set working directory",
				Value:   ".",
			},
			&cli.StringSliceFlag{
				Name:  "env-file",
				Usage: "set an alternate environment file",
				Value: []string{".env"},
			},
			&cli.StringFlag{
				Name:  "codex-path",
				Usage: "set codex executable path",
				Value: "codex",
			},
			&cli.StringFlag{
				Name:        "log-level",
				Usage:       "set log-level (\"error\", \"warn\", \"info\", \"debug\", \"trace\", \"off\", default: \"off\")",
				HideDefault: true,
				Value:       "off",
			},
		},
		Arguments: []cli.Argument{},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			// オプション引数の値を取得
			configPath := cmd.String("config")
			workdir := cmd.String("workdir")
			envFiles := cmd.StringSlice("env-file")
			codexPath := cmd.String("codex-path")
			logLevel := cmd.String("log-level")

			// OpenAI の API Key を取得
			apiKey, err := getAPIKey(ctx, appName, codexPath, envFiles)
			if err != nil {
				return err
			}

			// エージェントを定義したYAMLファイルを読み込み
			config, err := app.LoadConfig(configPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to load agent defined YAML file.\n")
				return fmt.Errorf("%w: %s", ErrInternal, err)
			}

			// エージェント名のチェック
			if cmd.Args().Len() == 0 {
				fmt.Fprintf(os.Stderr, "Please specify AGENT_NAME.\n")
				return fmt.Errorf("%w: %s", ErrUsage, err)
			}

			// アプリケーションをつくり、MCP Serverを実行
			app := app.New(
				config,
				codexPath,
				apiKey,
				subAgentMCPServerConfig(configPath, workdir, codexPath, apiKey),
				app.WithLogger(os.Stderr, logLevel),
			)
			agentName := cmd.Args().First()
			if err := app.RunMCPServer(agentName, workdir); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to start MCP server.\n")
				return fmt.Errorf("%w: %s", ErrInternal, err)
			}

			return nil
		},
	}
}
