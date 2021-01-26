package main

import (
	"fmt"
	"log"
	"sync/atomic"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/yomorun/yomo-source-mqtt-broker-starter/pkg/env"
)

var (
	brokerAddr = env.GetString("YOMO_SOURCE_MQTT_BROKER_ADDR", "tcp://0.0.0.0:1883")
	topic      = env.GetString("YOMO_SOURCE_MQTT_PUB_TOPIC", "NOISE")
	interval   = env.GetInt("YOMO_SOURCE_MQTT_PUB_INTERVAL", 100)
	counter    int64
)

func main() {
	options := mqtt.NewClientOptions().
		AddBroker(brokerAddr).
		SetUsername("admin").
		SetPassword("public")
	log.Println("Broker Addresses: ", options.Servers)
	options.SetClientID(fmt.Sprintf("yomo-source-pub-%d", time.Now().Unix()))
	options.SetConnectTimeout(time.Duration(0) * time.Second)
	options.SetAutoReconnect(true)
	options.SetKeepAlive(time.Duration(20) * time.Second)
	options.SetMaxReconnectInterval(time.Duration(5) * time.Second)
	options.OnConnect = func(client mqtt.Client) {
		fmt.Println("Connected")
	}
	options.OnConnectionLost = func(client mqtt.Client, err error) {
		fmt.Printf("Connect lost: %v", err)
	}
	options.SetOnConnectHandler(func(c mqtt.Client) {
		log.Printf("[client connect state] IsConnected:%v, IsConnectionOpen:%v", c.IsConnected(), c.IsConnectionOpen())
	})
	options.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		log.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	})

	client := mqtt.NewClient(options)
	for {
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			fmt.Printf("%v: Connect error:%v\n", time.Now().Format("2006-01-02 15:04:05"), token.Error())
			time.Sleep(time.Duration(1) * time.Second)
			continue
		}
		break
	}

	for {
		counter = atomic.AddInt64(&counter, 1)

		payload := fmt.Sprintf("{\"noise\":%v}", counter)
		pub(client, topic, payload)

		fmt.Printf("%v: Publish counter=%d, topic=%v, payload=%v\n", time.Now().Format("2006-01-02 15:04:05"), counter, topic, payload)
		time.Sleep(time.Duration(interval) * time.Millisecond)
	}
}

func pub(client mqtt.Client, topic string, payload interface{}) {
	for {
		if token := client.Publish(topic, 1, false, payload); token.Wait() && token.Error() != nil {
			log.Printf("yomo-source Publish error: %s \n", token.Error())
			time.Sleep(time.Duration(interval) * time.Millisecond)
			continue
		}
		break
	}
}
