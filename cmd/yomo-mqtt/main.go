package main

import (
	"github.com/spf13/cobra"
	"github.com/yomorun/yomo-source-mqtt-broker-starter/internal/cmd"
)

func main() {
	rootCmd := &cobra.Command{
		Use: "yomo-mqtt",
	}

	rootCmd.AddCommand(
		cmd.NewCmdBuild(),
		cmd.NewCmdRun(),
	)

	rootCmd.Execute()
}
