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
