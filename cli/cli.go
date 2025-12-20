package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/kurusugawa-computer/ace/agents"
	"github.com/kurusugawa-computer/ace/cli/credentials"
	"github.com/thamaji/codex-go"
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
func subAgentMCPServerConfig(configPath string, workdir string, codexPath string, apiKey string) func(subAgent *agents.SubAgent) (map[string]any, error) {
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
				"--codex-path",
				codexPath,
				subAgent.Name,
			},
			"startup_timeout_sec": 30,
			"tool_timeout_sec":    subAgent.TimeoutSec,
		}
		if apiKey != "" {
			config["env"] = map[string]any{
				"OPENAI_API_KEY": apiKey,
			}
		}

		return config, nil
	}
}

// OpenAI の API Key を取得
// 優先順位：Codex CLI のログイン状況 > 環境変数OPENAI_API_KEY > envfileオプションで指定された.envファイル > 設定ファイル
func getAPIKey(ctx context.Context, appName string, codexPath string, envFiles []string) (string, error) {
	// Codex CLI のログイン状況
	codexInstance := codex.New(codex.WithExecutablePath(codexPath))
	loggedIn, err := codexInstance.IsLoggedIn(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to check codex login status.\n")
		return "", fmt.Errorf("%w: %s", ErrInternal, err)
	}
	if loggedIn {
		return "", nil
	}

	// 環境変数 OPENAI_API_KEY
	// envfileオプションで指定された.envファイル
	_ = godotenv.Load(envFiles...)
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey != "" {
		return apiKey, nil
	}

	// 設定ファイル
	credentials, err := credentials.Load(appName)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			fmt.Fprintf(os.Stderr, "Failed to load credentials file.\n")
			return "", fmt.Errorf("%w: %s", ErrInternal, err)
		}

		fmt.Fprintf(os.Stderr, "The OpenAI API key is not set.\n")
		fmt.Fprintf(os.Stderr, "Please specify the environment variable OPENAI_API_KEY or run the `%s setup` command.\n", filepath.Base(os.Args[0]))
		return "", fmt.Errorf("%w: %s", ErrUsage, err)
	}

	return credentials.OpenAIAPIKey, nil
}
