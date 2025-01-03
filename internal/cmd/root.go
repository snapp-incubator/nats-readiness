package cmd

import (
	"log"
	"os"

	"github.com/snapp-incubator/nats-readiness/internal/cmd/serve"
	"github.com/spf13/cobra"
)

// ExitFailure status code.
const ExitFailure = 1

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	//nolint: exhaustruct
	root := &cobra.Command{
		Use:   "nats-readiness",
		Short: " NATS Readiness Gate ",
	}

	serve.Register(root)

	if err := root.Execute(); err != nil {
		log.Printf("failed to execute root command %s", err)
		os.Exit(ExitFailure)
	}
}
