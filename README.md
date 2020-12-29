# yomo-source-mqtt-borker-starter
Receive MQTT messages and convert them to the YoMo protocol for transmission to Serverless Service.



## 🚀 Getting Started

### Example (noise)

This example shows how to use the component reference method to make it easier to receive MQTT messages using starter and convert them to the YoMo protocol for transmission to the Zipper service.

#### 1. Init Project
```bash
go mod init source
go get github.com/yomorun/yomo-source-mqtt-borker-starter
```
#### 2. create app.go 
```text
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/yomorun/yomo/pkg/quic"

	"github.com/yomorun/yomo-codec-golang/pkg/codes"

	"github.com/yomorun/yomo-source-mqtt-borker-starter/pkg/starter"
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
			defer func() {
				if err := recover(); err != nil {
					log.Printf("handle error: %v", err)
				}
			}()
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
```

#### 3. run 
```bash
go run app.go
```

### Example (using cli: yomo-mqtt)

Running the application through the CLI

#### 1. build cli 
```bash
make build_cli
# create cli file: bin/yomo-mqtt
```
#### 2. create app.go
```text
package main

import (
	"encoding/json"
	"io"
	"log"

	"github.com/yomorun/yomo-codec-golang/pkg/codes"
)

func Handler(topic string, payload []byte, writer io.Writer) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("handle error: %v", err)
		}
	}()
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

	_, err = writer.Write(sendingBuf)
}
```

#### 3. run
```bash
./yomo-mqtt run -f app.go -p 1883 -z localhost:9999 -t NOISE
```
- -f Source function file (default is app.go)
- -p Port is the port number of MQTT host for Source function (default is 6262)
- -z Endpoint of ZipperAddr Server (default is localhost:4242)
- -t Topic of MQTT (default is NOISE)