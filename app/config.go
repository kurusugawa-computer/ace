package app

import (
	"os"

	"github.com/goccy/go-yaml"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/kurusugawa-computer/ace/agents"
)

type Config struct {
	// Codex CLI に与える config.toml
	// 詳細は https://github.com/openai/codex/blob/main/docs/config.md を参照。
	// ここでは、YAML ファイルに定義されているすべての AI エージェントに適用する Config を指定する。
	Config agents.CodexConfig `yaml:"config,omitempty"`

	// AI エージェントの定義
	// Key がエージェントの名前であり、CLI の実行時に指定する AGENT_NAME であり、
	// sub_agents で指定するサブエージェント名でもある。
	Agents map[string]*AgentConfig `yaml:"agents,omitempty"`
}

type AgentConfig struct {
	// AI エージェントの名前
	// YAML ファイルには記載しない。
	Name string `yaml:"-"`

	// AI エージェントの説明
	// AI エージェントをサブエージェントとして呼び出すとき、
	// もしくは MCP Server として利用するとき、
	// MCP Client に提供するツールの description として利用される。
	Description string `yaml:"description"`

	// AI エージェントに対する基本的な指示
	Instruction string `yaml:"instruction"`

	// AI エージェントに対するプロンプトのテンプレート
	// AI エージェントの実行時に input_schema で定義した Key-Value が展開される。
	// テンプレートの書式は golang 標準ライブラリの text/template（https://pkg.go.dev/text/template）のもの。
	// prompt_template が `「{{.question}}」` であり、input_schema が `{"question":{"type":"string"}}` であるとき
	// ユーザーの入力が `question=こんにちは` ならば、プロンプトは `「こんにちは」` となる。
	PromptTemplate string `yaml:"prompt_template"`

	// 入力スキーマ
	// ユーザー、もしくは MCP Client からの入力データの形式を JSON Schema 形式で定義する。
	// AI が参照するので、description を丁寧に書くことを推奨する。
	InputSchema map[string]*jsonschema.Schema `yaml:"input_schema"`

	// 出力スキーマ
	// AI エージェントの出力データの形式を JSON Schema 形式で定義する。
	// AI が参照するので、description を丁寧に書くことを推奨する。
	OutputSchema map[string]*jsonschema.Schema `yaml:"output_schema"`

	// Codex がユーザーの承認を求めるタイミング
	// https://github.com/openai/codex/blob/main/docs/config.md#approval_policy を参照。
	// デフォルト値は never
	ApprovalPolicy string `yaml:"approval_policy"` // untrusted, on-failure, on-request, never

	// Codex のサンドボックスモード
	// https://github.com/openai/codex/blob/main/docs/config.md#sandbox_mode を参照。
	// デフォルト値は read-only
	Sandbox string `yaml:"sandbox"` // read-only, workspace-write, danger-full-access

	// AI エージェントをサブエージェントとして実行したとき、タイムアウトする秒数
	TimeoutSec int `yaml:"timeout_sec,omitempty"` // default: 1800

	// 利用する MCP Server の定義
	// https://github.com/openai/codex/blob/main/docs/config.md#mcp_servers を参照。
	MCPServers map[string]MCPServerConfig `yaml:"mcp_servers"`

	// 利用するサブエージェントのエージェント名のリスト
	// サブエージェントを指定すると、同じ YAML ファイルに定義されている別の AI エージェントを
	// AI エージェントの実行中にツールとして呼び出して利用できる。
	SubAgents []string `yaml:"sub_agents,omitempty"`

	// Codex CLI に与える config.toml
	// 詳細は https://github.com/openai/codex/blob/main/docs/config.md を参照。
	// ここでは、この AI エージェントにのみ適用する Config を指定する。
	Config agents.CodexConfig `yaml:"config,omitempty"`
}

type MCPServerConfig map[string]any

func LoadConfig(path string) (*Config, error) {
	f, err := os.OpenFile(path, 0, os.FileMode(os.O_RDONLY))
	if err != nil {
		return nil, err
	}

	dec := yaml.NewDecoder(f, yaml.UseJSONUnmarshaler())
	config := Config{}
	err = dec.Decode(&config)
	f.Close()
	if err != nil {
		return nil, err
	}

	// エージェントConfigのkeyをエージェントのNameとしてセット
	for name, agentConfig := range config.Agents {
		agentConfig.Name = name
	}

	return &config, nil
}
