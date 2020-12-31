package source

import (
	"context"
	"io"
	"log"
	"plugin"

	"github.com/yomorun/yomo-source-mqtt-broker-starter/pkg/starter"
	yquic "github.com/yomorun/yomo/pkg/quic"
)

func Run(addr string, handler *MQTTServerHandler) error {
	client, err := yquic.NewClient(addr)
	if err != nil {
		log.Printf("NewClient addr=%s error:%s", addr, err.Error())
		return err
	}

	stream, err := client.CreateStream(context.Background())
	if err != nil {
		log.Printf("CreateStream addr=%s error:%s", addr, err.Error())
		return err
	}

	starter.NewBrokerSimply(handler.Endpoint, handler.Topic).
		Run(func(topic string, payload []byte) {
			payloadHandlerFunc, ok := handler.Handler.(func(string, []byte, io.Writer))
			if !ok {
				log.Fatalln("payloadHandlerFunc error")
			}

			payloadHandlerFunc(topic, payload, stream)
		})

	return nil
}

type MQTTServerHandler struct {
	Handler  plugin.Symbol
	Endpoint string
	Topic    string
}
