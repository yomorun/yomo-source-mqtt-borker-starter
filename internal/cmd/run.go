package cmd

import (
	"fmt"
	"log"

	"github.com/yomorun/yomo-source-mqtt-broker-starter/internal/source"

	"github.com/spf13/cobra"
)

type RunOptions struct {
	baseOptions

	ZipperAddr string
	Port       int
	Topic      string
}

func NewCmdRun() *cobra.Command {
	var opts = &RunOptions{}

	var cmd = &cobra.Command{
		Use:   "run",
		Short: "Run a YoMo Source Function",
		Long:  "Run a YoMo Source Function.",
		Run: func(cmd *cobra.Command, args []string) {
			soHandler, err := buildAndLoadHandler(&opts.baseOptions, args)
			if err != nil {
				return
			}

			// serve the Source app
			endpoint := fmt.Sprintf("0.0.0.0:%d", opts.Port)
			handler := &source.MQTTServerHandler{
				Handler:  soHandler,
				Endpoint: endpoint,
				Topic:    opts.Topic,
			}

			err = source.Run(opts.ZipperAddr, handler)
			if err != nil {
				log.Print("Run the serverless failure with err: ", err)
			}
		},
	}

	cmd.Flags().StringVarP(&opts.Filename, "file-name", "f", "app.go", "Source function file (default is app.go)")
	cmd.Flags().IntVarP(&opts.Port, "port", "p", 1883, "Port is the port number of MQTT host for Source function (default is 6262)")
	cmd.Flags().StringVarP(&opts.ZipperAddr, "zipper-addr", "z", "localhost:9999", "Endpoint of ZipperAddr Server (default is localhost:4242)")
	cmd.Flags().StringVarP(&opts.Topic, "topic", "t", "NOISE", "Topic of MQTT (default is NOISE)")

	return cmd
}
