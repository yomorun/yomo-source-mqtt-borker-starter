package main

import (
	"encoding/json"
	"io"
	"log"

	y3 "github.com/yomorun/y3-codec-golang"
)

func Handler(topic string, payload []byte, writer io.Writer) {
	log.Printf("topic=%v, payload=%v\n", topic, string(payload))

	// get data from MQTT
	var raw map[string]int32
	err := json.Unmarshal(payload, &raw)
	if err != nil {
		log.Printf("Unmarshal payload error:%v", err)
	}

	// generate YoMo-Codec format
	data := float32(raw["noise"])
	codec := y3.NewCodec(0x10)
	sendingBuf, _ := codec.Marshal(data)
	log.Printf("sendingBuf=%#x\n", sendingBuf)

	_, err = writer.Write(sendingBuf)
}
