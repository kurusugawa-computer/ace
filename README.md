# ace - Agent Command Executor

<p align="center">
  <img width="400" alt="ace" src="./ace_logo.png" /><br/>
ACE は、AI エージェントを YAML ファイルで簡潔に定義し、CLI から実行できるツールです。<br/>
サブエージェントによるマルチエージェント構成や MCP Server の利用にも対応しています。
</p>

## Description

ACE は、AI エージェントの定義・実行を効率化するためのコマンドラインツールです。  
エージェントの動作ロジック、入力・出力仕様、モデル設定などを YAML で記述するだけで、柔軟なエージェント実行環境を構築できます。

AI エージェントの定義がひとつの YAML ファイルで完結するので、YAML ファイルさえ共有してしまえば、構築した AI エージェントをチームで共有することも簡単です。

## Installation

### 必要環境

ACE は内部的に OpenAI の Codex CLI を利用しています。  
`codex` コマンドにパスが通っている必要があります。

また、必要に応じて MCP Server の実行環境（npm や uv など）を用意してください。

### インストール手順

1. GitHub の Release ページから最新のアーカイブをダウンロードする  
2. 任意のディレクトリに展開し、実行バイナリへパスを通す

## Setup

ACE の実行には **OpenAI API Key** が必要です。以下のいずれかの方法で設定してください。

- 環境変数 `OPENAI_API_KEY` を指定する
- `.env` ファイルに API Key を記載する
- `ace setup` コマンドを用いて対話的に設定する

## Usage

ACE では、AI エージェントを YAML ファイルで定義します。  
YAML ファイルのスキーマは [app/config.go#L11](https://github.com/kurusugawa-computer/ace/blob/main/app/config.go#L11) を参照してください。

次の YAML ファイルは、ユーザーの質問にたいして元気よく回答する AI エージェントです。  
ユーザーが天気について質問すると、この AI エージェントは weather サブエージェントを実行して天気情報を取得します。  
weather エージェントは MCP Server として Web 検索を行うツールを与えられています。

```yaml
config:
  model_provider: openai
  model: gpt-5-mini
  model_reasoning_effort: low
  model_verbosity: low

agents:
  root:
    description: |
      ユーザーの質問に回答します。
    instruction: |
      ユーザーの質問に対して元気よく回答すること。
    prompt_template: |
      {{.question}}
    input_schema:
      question:
        type: string
    output_schema:
      answer:
        type: string
    approval_policy: untrusted
    sub_agents:
      - weather

  weather:
    description: |
      天気情報を調べます。
    instruction: |
      以下のツールを使用して、天気情報を検索してユーザーの質問に回答しなさい。
      - search
      - fetch_content
    prompt_template: |
      {{.location}} の {{.time}} の天気について Web 検索して回答しなさい。
    input_schema:
      location:
        type: string
        description: 地域
      time:
        type: string
        description: 日時
    output_schema:
      result:
        type: string
        description: 天気情報
    mcp_servers:
      duckduckgo:
        command: uvx
        args: ["duckduckgo-mcp-server"]
        enabled_tools:
          - search
          - fetch_content
```

この YAML ファイルが `examples/simple.yaml` に保存されているとき、定義した `root` エージェントを実行するには、次のコマンドを実行します。

```bash
ace exec -c example/simple.yaml root question=明日の名古屋の天気は？
```

実行結果は以下のような JSON 形式で出力されます。

```json
{
  "answer": "明日の名古屋は晴のち雪の予報です！予想最高気温は13℃、最低気温は9℃。降水確率は午前〜昼が約30%、夕方以降は50%と上がります。風は北西で最大7m/sほどです。出かける際は念のため防寒と雨具をお持ちくださいね！"
}
```

コマンドの詳細は `ace --help` や `ace exec --help` を参照してください。

### MCP Server としての利用

mcp-server サブコマンドを実行すると、ACE を MCP Server（STDIO 形式）として起動できます。  
次のコマンドでは、YAML ファイルに定義されている `root` エージェントを、`root` ツールとして MCP Client へ提供します。

```bash
ace mcp-server -c example/simple.yaml root
```

### プログラムからの利用

以下のバインディングライブラリを利用できます。

| プログラミング言語 | ライブラリ                                              |
| ------------------ | ------------------------------------------------------- |
| Golang             | [ace-go](https://github.com/kurusugawa-computer/ace-go) |
