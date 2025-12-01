package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/kurusugawa-computer/ace/cli"
)

var appName = "ace"
var version = "v0.0.0"
var title = "Agent Command Executor"

func main() {
	app := cli.New(
		appName,
		version,
		title,
	)

	ctx := context.Background()

	if err := app.Run(ctx, os.Args); err != nil {
		switch {
		case errors.Is(err, cli.ErrUsage):
			// ErrUsage の場合はすでにエラーメッセージが出ているはずなので、何もしない

		case errors.Is(err, cli.ErrInternal):
			fmt.Fprintf(os.Stderr, "%s\n", err)

		default:
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		}

		os.Exit(1)
	}
}
