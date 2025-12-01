package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/kurusugawa-computer/ace/app"
	"github.com/kurusugawa-computer/ace/cli/credentials"
	"github.com/urfave/cli/v3"
)

var _ subCommand = exec

func exec(appName string, version string) *cli.Command {
	return &cli.Command{
		Name:    "exec",
		Aliases: []string{},
		Usage: `Execute an AI agent defined in a YAML file.
KEY=VALUE pairs can be used to fill prompt_template variables.`,
		ArgsUsage: "AGENT_NAME [KEY=VALUE...]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "set YAML file where the AI ​​agent is defined",
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
			logLevel := cmd.String("log-level")

			// OpenAI の API Key を取得
			// 優先順位：環境変数OPENAI_API_KEY > envfileオプションで指定された.envファイル > 設定ファイル
			_ = godotenv.Load(envFiles...)
			apiKey := os.Getenv("OPENAI_API_KEY")
			if apiKey == "" {
				credentials, err := credentials.Load(appName)
				if err != nil {
					if !errors.Is(err, os.ErrNotExist) {
						fmt.Fprintf(os.Stderr, "Failed to load credentials file.\n")
						return fmt.Errorf("%w: %s", ErrInternal, err)
					}
					fmt.Fprintf(os.Stderr, "The OpenAI API key is not set.\n")
					fmt.Fprintf(os.Stderr, "Please specify the environment variable OPENAI_API_KEY or run the `%s setup` command.\n", filepath.Base(os.Args[0]))
					return fmt.Errorf("%w: %s", ErrUsage, err)
				}
				apiKey = credentials.OpenAIAPIKey
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

			// コマンドライン引数をパースしてAIエージェントの入力データをつくる
			input := map[string]any{}
			for _, argument := range cmd.Args().Tail() {
				segments := strings.SplitN(argument, "=", 2)
				if len(segments) != 2 {
					fmt.Fprintf(os.Stderr, "Invalid KEY=VALUE argument: %s\n", argument)
					return fmt.Errorf("%w: %s", ErrUsage, err)
				}
				key := segments[0]
				value := segments[1]
				input[key] = value
			}

			// アプリケーションをつくり、AIエージェントを実行
			app := app.New(
				config,
				apiKey,
				subAgentMCPServerConfig(configPath, workdir, apiKey),
				app.WithLogger(os.Stderr, logLevel),
			)
			agentName := cmd.Args().First()
			output, err := app.RunAgent(agentName, workdir, input)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to start AI agent\n")
				return fmt.Errorf("%w: %s", ErrInternal, err)
			}

			// AIエージェントの実行結果をJSON形式で出力
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			_ = enc.Encode(output)

			return nil
		},
	}
}
