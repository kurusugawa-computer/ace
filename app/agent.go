package app

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/kurusugawa-computer/ace/agents"
)

// エージェントを実行する。
// arguments: ["KEY=VALUE", "KEY=VALUE"...]
func (app *App) RunAgent(agentName string, workdir string, arguments []string) (any, error) {
	// エージェントのビルド
	agent, err := app.buildAgent(agentName)
	if err != nil {
		return nil, err
	}

	// ["KEY=VALUE", "KEY=VALUE"...] 形式の arguments を map[string]any にパースする
	// ただし、この段階では value は string のまま。(KEYの階層構造のみをパース)
	argumentsMap, err := parseArguments(arguments)
	if err != nil {
		return nil, err
	}

	// input_schema の定義に従って、arguments をパースする
	input := map[string]any{}
	for propName, propSchema := range agent.InputSchema.Properties {
		value, err := applyJSONSchema(propName, argumentsMap[propName], propSchema)
		if err != nil {
			return nil, err
		}
		input[propName] = value
	}

	// vars の値を展開する
	if app.config.Vars != nil {
		for key, value := range app.config.Vars {
			if _, ok := input[key]; ok {
				continue
			}
			input[key] = value
		}
	}

	// エージェントの実行
	output, err := agent.Run(
		workdir,
		input,
		&agents.RunConfig{
			APIKey:                  app.apiKey,
			SubagentMCPServerConfig: app.subAgentMCPServerConfig,
			LogLevel:                app.logLevel,
			LogWriter:               app.logWriter,
		},
	)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func parseArguments(arguments []string) (map[string]any, error) {
	input := map[string]any{}
	for _, argument := range arguments {
		segments := strings.SplitN(argument, "=", 2)
		if len(segments) != 2 {
			return nil, fmt.Errorf("argument is not in KEY=VALUE format: %s", argument)
		}

		keys := strings.Split(segments[0], ".")
		value := segments[1]

		currentMap := input
		for i := 0; i < len(keys)-1; i++ {
			key := keys[i]

			if value, ok := currentMap[key]; ok {
				child, ok := value.(map[string]any)
				if !ok {
					return nil, fmt.Errorf("conflicting input values specified in arguments: %s", strings.Join(keys[0:i], "."))
				}
				currentMap = child
			} else {
				child := map[string]any{}
				currentMap[key] = child
				currentMap = child
			}
		}

		key := keys[len(keys)-1]
		currentMap[key] = value
	}

	return input, nil
}

func applyJSONSchema(key string, value any, schema *jsonschema.Schema) (any, error) {
	switch schema.Type {
	default:
		return value, nil

	case "string":
		if value == nil {
			if schema.Default == nil {
				return nil, fmt.Errorf("missing required field specified in input_schema: %s", key)
			}

			var str string
			if err := json.Unmarshal(schema.Default, &str); err != nil {
				return nil, err
			}

			return str, nil
		}

		value, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("specified value is incompatible with the input_schema definition: %s", key)
		}

		return value, nil

	case "integer":
		if value == nil {
			if schema.Default == nil {
				return nil, fmt.Errorf("missing required field specified in input_schema: %s", key)
			}

			var integer int64
			if err := json.Unmarshal(schema.Default, &integer); err != nil {
				return nil, err
			}

			return integer, nil
		}

		value, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("specified value is incompatible with the input_schema definition: %s", key)
		}

		integer, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, err
		}

		return integer, nil

	case "number":
		if value == nil {
			if schema.Default == nil {
				return nil, fmt.Errorf("missing required field specified in input_schema: %s", key)
			}

			var number float64
			if err := json.Unmarshal(schema.Default, &number); err != nil {
				return nil, err
			}

			return number, nil
		}

		value, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("specified value is incompatible with the input_schema definition: %s", key)
		}

		number, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, err
		}

		return number, nil

	case "boolean":
		if value == nil {
			if schema.Default == nil {
				return nil, fmt.Errorf("missing required field specified in input_schema: %s", key)
			}

			var boolean float64
			if err := json.Unmarshal(schema.Default, &boolean); err != nil {
				return nil, err
			}

			return boolean, nil
		}

		value, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("specified value is incompatible with the input_schema definition: %s", key)
		}

		boolean, err := strconv.ParseBool(value)
		if err != nil {
			return nil, err
		}

		return boolean, nil

	case "null":
		if value == nil {
			return nil, nil
		}

		value, ok := value.(string)
		if !ok || value != "null" {
			return nil, fmt.Errorf("specified value is incompatible with the input_schema definition: %s", key)
		}

		return nil, nil

	case "array":
		if value == nil {
			if schema.Default == nil {
				return nil, fmt.Errorf("missing required field specified in input_schema: %s", key)
			}

			var array []any
			if err := json.Unmarshal(schema.Default, &array); err != nil {
				return nil, err
			}

			return array, nil
		}

		values, ok := value.([]string)
		if !ok {
			return nil, fmt.Errorf("specified value is incompatible with the input_schema definition: %s", key)
		}

		schemas := make([]*jsonschema.Schema, 0, len(values))
		schemas = append(schemas, schema.PrefixItems...)
		itemsSchema := schema.UnevaluatedItems
		if schema.Contains != nil {
			itemsSchema = schema.Contains
		}
		if schema.Items != nil {
			itemsSchema = schema.Items
		}
		if itemsSchema != nil {
			for i := len(schemas); i < len(values); i++ {
				schemas = append(schemas, itemsSchema)
			}
		}

		array := make([]any, 0, len(values))
		for i := 0; i < len(values); i++ {
			value, err := applyJSONSchema(fmt.Sprintf("%s[%d]", key, i), values[i], schemas[i])
			if err != nil {
				return nil, err
			}

			array = append(array, value)
		}

		return array, nil

	case "object":
		if value == nil {
			if schema.Default == nil {
				return nil, fmt.Errorf("missing required field specified in input_schema: %s", key)
			}

			var object map[string]any
			if err := json.Unmarshal(schema.Default, &object); err != nil {
				return nil, err
			}

			return object, nil
		}

		values, ok := value.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("specified value is incompatible with the input_schema definition: %s", key)
		}

		object := map[string]any{}
		for propName, propSchema := range schema.Properties {
			value, err := applyJSONSchema(fmt.Sprintf("%s.%s", key, propName), values[propName], propSchema)
			if err != nil {
				return nil, err
			}

			object[propName] = value
		}

		return object, nil
	}
}
