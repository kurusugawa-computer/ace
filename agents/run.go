package agents

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/thamaji/codex-go"
)

type RunConfig struct {
	APIKey                  string
	SubagentMCPServerConfig func(subAgent *SubAgent) (map[string]any, error)
	LogLevel                string // error, warn, info, debug, trace, off
	LogWriter               io.Writer
}

func (agent *Agent) Run(workdir string, input map[string]any, config *RunConfig) (any, error) {
	// 作業ディレクトリの絶対パスを取得
	workdirAbsPath, err := filepath.Abs(workdir)
	if err != nil {
		return nil, err
	}

	// Codex の Config を構築
	codexConfig := agent.Config.Clone()
	for _, subAgent := range agent.SubAgents {
		mcpServerConfig, err := config.SubagentMCPServerConfig(subAgent)
		if err != nil {
			return nil, err
		}

		codexConfig["mcp_servers."+subAgent.Name] = mcpServerConfig
	}

	// AGENTS.md が存在しても見に行かないように制限
	codexConfig["project_doc_max_bytes"] = 0

	// プロンプトの構築
	prompt := &strings.Builder{}
	if err := agent.PromptTemplate.Execute(prompt, input); err != nil {
		return nil, err
	}

	// プロンプトに出力形式の指定を追加
	outputSchemaJSON, err := agent.OutputSchema.MarshalJSON()
	if err != nil {
		return nil, err
	}
	fmt.Fprintln(prompt, "")
	fmt.Fprintln(prompt, "なお、回答の出力形式は以下の JSON Schema に厳格に従うこと。")
	fmt.Fprintln(prompt, string(outputSchemaJSON))

	// Codex を実行して回答を取得
	var options []codex.CodexOption
	if agent.codexExecutablePath != "" {
		options = append(options, codex.WithExecutablePath(agent.codexExecutablePath))
	}
	if config.LogWriter != nil && config.LogLevel != "off" {
		options = append(options, codex.WithLogger(config.LogWriter, config.LogLevel))
	}
	codexInstance := codex.New(options...)
	ctx := context.Background()

	loggedIn, err := codexInstance.IsLoggedIn(ctx)
	if err != nil {
		return nil, err
	}
	if !loggedIn {
		if err := codexInstance.Login(ctx, config.APIKey); err != nil {
			return nil, err
		}
	}

	answer, err := codexInstance.Invoke(
		ctx,
		prompt.String(),
		codex.WithDeveloperInstructions(agent.Instruction),
		codex.WithCwd(workdirAbsPath),
		codex.WithApprovalPolicy(agent.ApprovalPolicy),
		codex.WithSandbox(agent.Sandbox),
		codex.WithConfig(codexConfig),
	)
	if err != nil {
		return nil, err
	}
	answer = strings.TrimSpace(answer)

	// 回答が出力形式に従っているかチェック
	var output any
	if err := json.Unmarshal([]byte(answer), &output); err == nil {
		resolved, err := agent.OutputSchema.Resolve(nil)
		if err != nil {
			return nil, err
		}
		if resolved.Validate(output) == nil {
			// 出力形式に従っていたら、そのまま返す
			return output, nil
		}
	}

	// 回答の内容を AI で出力形式に合わせて整形する
	client := openai.NewClient(option.WithAPIKey(config.APIKey))
	chatCompletion, err := client.Chat.Completions.New(
		ctx,
		openai.ChatCompletionNewParams{
			Model:       openai.ChatModelGPT5Nano,
			Temperature: openai.Float(1),
			ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
				OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
					JSONSchema: openai.ResponseFormatJSONSchemaJSONSchemaParam{
						Name:        "parse",
						Description: openai.String("ユーザーの入力した文章を解釈します。"),
						Schema:      agent.OutputSchema,
						Strict:      openai.Bool(true),
					},
				},
			},
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage("ユーザーの入力した文章をJSON Schemaに従って出力してください。"),
				openai.UserMessage(answer),
			},
		},
	)
	if err != nil {
		return nil, err
	}
	if chatCompletion == nil || len(chatCompletion.Choices) == 0 {
		return nil, errors.New("invalid format, openai chat completions response")
	}
	if err := json.Unmarshal([]byte(chatCompletion.Choices[0].Message.Content), &output); err != nil {
		return nil, errors.New("invalid format, openai chat completions response")
	}

	return output, nil
}
