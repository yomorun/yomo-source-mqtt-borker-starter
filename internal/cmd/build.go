package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

type BuildOptions struct {
	baseOptions
}

func NewCmdBuild() *cobra.Command {
	var opts = &BuildOptions{}

	var cmd = &cobra.Command{
		Use:   "build",
		Short: "Build the YoMo Source Function",
		Long:  "Build the YoMo Source Function as .so file",
		Run: func(cmd *cobra.Command, args []string) {
			_, err := buildSourceFile(&opts.baseOptions, args)
			if err != nil {
				return
			}
			log.Print("✅ Build the source file successfully.")
		},
	}

	cmd.Flags().StringVarP(&opts.Filename, "file-name", "f", "app.go", "Source function file (default is app.go)")

	return cmd
}
