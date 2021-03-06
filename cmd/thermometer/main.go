package main

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/yomorun/y3-codec-golang"
	"github.com/yomorun/yomo-source-mqtt-broker-starter/pkg/utils"

	"github.com/yomorun/yomo-source-mqtt-broker-starter/pkg/env"

	"github.com/yomorun/yomo/pkg/quic"

	"github.com/yomorun/yomo-source-mqtt-broker-starter/pkg/starter"
)

var (
	zipperAddr = env.GetString("YOMO_SOURCE_MQTT_ZIPPER_ADDR", "localhost:9999")
	brokerAddr = env.GetString("YOMO_SOURCE_MQTT_BROKER_ADDR", "0.0.0.0:1883")
)

type ThermometerData struct {
	Temperature float32 `y3:"0x11" json:"tem"`
	Humidity    float32 `y3:"0x12" json:"hum"`
}

func main() {
	var (
		stream = createStream()
		mutex  sync.Mutex
	)

	starter.NewBrokerSimply(brokerAddr, "thermometer").
		Run(func(topic string, payload []byte) {
			log.Printf("receive: topic=%v, payload=%v\n", topic, string(payload))

			// get data from MQTT
			var data ThermometerData
			err := json.Unmarshal(payload, &data)
			if err != nil {
				log.Printf("Unmarshal payload error:%v", err)
			}

			// generate y3-codec format
			sendingBuf, _ := y3.NewCodec(0x10).Marshal(data)

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

			log.Printf("write: sendingBuf=%v\n", utils.FormatBytes(sendingBuf))
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
