package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	cre "github.com/kurusugawa-computer/ace/cli/credentials"
	"github.com/urfave/cli/v3"
	"golang.org/x/term"
)

var _ subCommand = setup

func setup(appName string, version string) *cli.Command {
	return &cli.Command{
		Name:      "setup",
		Aliases:   []string{},
		Usage:     "Register your OpenAI API Key and setup " + appName + ".",
		ArgsUsage: " ",
		Flags:     []cli.Flag{},
		Arguments: []cli.Argument{},
		Action: func(aContext context.Context, aCommand *cli.Command) error {
			// OpenAI の API Key を入力してもらう
			openAIAPIKey, err := ReadPassword("input your OpenAI API key")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to read OpenAI API key.\n")
				return fmt.Errorf("%w: %s", ErrInternal, err)
			}

			// OpenAI の API Key を保存
			credentials := &cre.Credentials{
				OpenAIAPIKey: strings.TrimSpace(openAIAPIKey),
			}
			if err := cre.Save(appName, credentials); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to save credentials.\n")
				return fmt.Errorf("%w: %s", ErrInternal, err)
			}

			fmt.Printf("\n")
			fmt.Printf("Setup successful.\n")

			return nil
		},
	}
}

func ReadPassword(aMessage string) (string, error) {
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return "", errors.New("stdin is not a terminal")
	}

	fmt.Fprint(os.Stdout, aMessage+": ")
	tAnswer, tError := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Fprintln(os.Stdout)
	if tError != nil {
		return "", tError
	}

	return string(tAnswer), nil
}
