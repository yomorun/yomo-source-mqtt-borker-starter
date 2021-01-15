package main

import (
	"context"
	"encoding/json"
	"github.com/yomorun/yomo-source-mqtt-broker-starter/pkg/env"
	"log"
	"sync"

	"github.com/yomorun/yomo/pkg/quic"

	"github.com/yomorun/y3-codec-golang/pkg/codes"

	"github.com/yomorun/yomo-source-mqtt-broker-starter/pkg/starter"
)

var (
	zipperAddr = env.GetString("YOMO_SOURCE_MQTT_ZIPPER_ADDR", "localhost:9999")
	brokerAddr = env.GetString("YOMO_SOURCE_MQTT_BROKER_ADDR", "0.0.0.0:1883")
)

func main() {
	var (
		stream = createStream()
		mutex  sync.Mutex
	)

	starter.NewBrokerSimply(brokerAddr, "NOISE").
		Run(func(topic string, payload []byte) {
			log.Printf("receive: topic=%v, payload=%v\n", topic, string(payload))

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

			mutex.Lock()
			_, err = stream.Write(sendingBuf)
			if err != nil {
				log.Printf("stream.Write error: %v, sendingBuf=%#x\n", err, sendingBuf)
				err = stream.Close()
				if err != nil {
					log.Printf("stream.Close error: %v\n", err)
				}
				stream = createStream()
			}
			mutex.Unlock()

			log.Printf("write: sendingBuf=%#x\n", sendingBuf)
		})
}

func createStream() quic.Stream {
	var (
		err    error
		client quic.Client
		stream quic.Stream
	)

	for {
		client, err = quic.NewClient(zipperAddr)
		if err != nil {
			log.Printf("NewClient error: %v, addr=%v\n", err, zipperAddr)
			continue
		}
		break
	}

	for {
		stream, err = client.CreateStream(context.Background())
		if err != nil {
			log.Printf("CreateStream error: %v\n", err)
			continue
		}
		break
	}

	return stream
}
