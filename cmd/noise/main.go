package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/yomorun/yomo/pkg/quic"

	"github.com/yomorun/yomo-codec-golang/pkg/codes"

	"github.com/yomorun/yomo-source-mqtt-broker-starter/pkg/starter"
)

func main() {
	client, err := quic.NewClient("localhost:9999")
	if err != nil {
		panic(fmt.Errorf("NewClient error:%s", err.Error()))
	}

	stream, err := client.CreateStream(context.Background())
	if err != nil {
		panic(fmt.Errorf("CreateStream error:%s", err.Error()))
	}

	starter.NewBrokerSimply("localhost:1883", "NOISE").
		Run(func(topic string, payload []byte) {
			log.Printf("topic=%v, payload=%v\n", topic, string(payload))

			// get data from MQTT
			var raw map[string]int32
			err := json.Unmarshal(payload, &raw)
			if err != nil {
				log.Printf("Unmarshal payload error:%v", err)
			}

			// generate YoMo-Codec format
			data := float32(raw["noise"])
			proto := codes.NewProtoCodec(0x10)
			sendingBuf, _ := proto.Marshal(data)
			log.Printf("sendingBuf=%#x\n", sendingBuf)

			_, err = stream.Write(sendingBuf)
		})
}
