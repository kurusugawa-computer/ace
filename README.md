# ace - Agent Command Executor

<p align="center">
  <img width="400" alt="ace" src="./ace_logo.png" /><br/>
ACE は、AI エージェントを YAML ファイルで簡潔に定義し、CLI から実行できるツールです。<br/>
サブエージェントによるマルチエージェント構成や MCP Server の利用にも対応しています。
</p>

## Description

ACE は、AI エージェントの定義・実行を効率化するためのコマンドラインツールです。  
エージェントの動作ロジック、入力・出力仕様、モデル設定などを YAML で記述するだけで、柔軟なエージェント実行環境を構築できます。

## Installation

1. GitHub の Release ページから最新のアーカイブをダウンロードする  
2. 任意のディレクトリに展開し、実行バイナリへパスを通す

## Setup

ACE の実行には **OpenAI API Key** が必要です。以下のいずれかの方法で設定してください。

- 環境変数 `OPENAI_API_KEY` を指定する
- `.env` ファイルに API Key を記載する
- `ace setup` コマンドを用いて対話的に設定する

## Usage

### AI エージェント定義ファイル（YAML）

ACE では、YAML 形式でエージェントを定義します。  
以下はサンプルの定義です。

```yaml
agents:
  root:
    description: テスト
    instruction: 元気よく回答すること
    prompt_template: |
      {{.question}}
    input_schema:
      question:
        type: string
    output_schema:
      answer:
        type: string
    approval_policy: never
    sandbox: read-only
    config:
      model: gpt-5-nano
      model_reasoning_effort: low
      model_verbosity: high
      tools:
        web_search: false
    sub_agents:
      - weather

  weather:
    description: 天気情報を調べます
    instruction: 何を聞かれても「嵐」とだけ回答すること
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
    approval_policy: never
    sandbox: read-only
    config:
      model: gpt-5-nano
      model_reasoning_effort: low
      model_verbosity: low
      tools:
        web_search: false
```

#### ポイント

1. **agents**  
   エージェント名をキーとする辞書形式で定義します。  
   キーがエージェント名、値がその設定です。

2. **instruction / prompt_template**  
   エージェントの振る舞いを決定する主要要素です。  
   input_schema で定義した値は prompt_template 内で {{.key}} として参照できます。  
   例：`ace exec root question="今日の名古屋の天気は？"`  
   この場合、{{.question}} には "今日の名古屋の天気は？" が展開されます。  

3. **input_schema / output_schema**  
   JSON Schema を YAML で記述した形式で指定します。  
   トップレベルにはオブジェクトのみを使用できます（配列・文字列・数値の直接指定は不可）。

4. **description**  
   サブエージェントとして呼び出される場合、親エージェントに対する tool description として利用されます。  
    用途が明確にわかる形で記述してください。

5. **approval_policy / sandbox**  
   Codex の approval_policy、sandbox の値を指定します。  
   [Codex のドキュメント](https://github.com/openai/codex/blob/main/docs/config.md) を参照してください。

6. **config**  
   Codex の config.toml の YAML 表現です。  
   [Codex のドキュメント](https://github.com/openai/codex/blob/main/docs/config.md) を参照してください。

### 実行例

次の例は `agent.yaml` に定義されている `root` エージェントを、`question=今日の名古屋の天気は？` で実行します。

```bash
ace exec -c agent.yaml root question=今日の名古屋の天気は？
```

実行結果は以下のような JSON 形式で出力されます。（output_schema で指定できます。）

```json
{
  "answer": "名古屋の明日の天気は嵐の見込みです。強風と大雨に備え、傘や雨具、防風対策を準備してください。最新の天気情報をこまめにご確認ください。"
}
```

詳細は `ace --help` や `ace exec --help` を参照してください。

### MCP Server としての利用

次のように呼び出すと、ace を MCP Server（STDIO 形式）として起動できます。  
この場合、`root` エージェントが `root` ツールとして MCP Client へ提供されます。

```bash
ace mcp-server -c agent.yaml root
```
