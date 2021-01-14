package starter

import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"time"

	"github.com/fhmq/hmq/broker"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	client mqtt.Client
)

type Broker struct {
	conf *BrokerConf
}

func (b Broker) Run(handle func(topic string, payload []byte)) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		mqttRun(b.conf.Addr, &wg)
	}()
	wg.Wait()

	b.Sub(handle)
}

func (b Broker) Sub(handle func(topic string, payload []byte)) {
	var count = 0
	for {
		b.connSub(b.conf, func(client mqtt.Client, msg mqtt.Message) {
			go handle(msg.Topic(), msg.Payload())
		})
		for {
			if !client.IsConnected() || !client.IsConnectionOpen() {
				count = count + 1
				if count >= b.conf.FailureTimes {
					count = 0
					client.Disconnect(500)
					break
				}
			}
			time.Sleep(time.Duration(5) * time.Second)
		}
	}
}

func NewBrokerSimply(addr string, topic string) Broker {
	return NewBroker(&BrokerConf{
		Addr:   addr,
		Topics: []string{topic},
	})
}

func NewBroker(conf *BrokerConf) Broker {
	broker := Broker{}

	if conf == nil {
		broker.conf = DefaultConfig
	} else {
		broker.conf = conf
		if conf.ConnectTimeout < 0 {
			conf.ConnectTimeout = DefaultConfig.ConnectTimeout
		}
		if len(conf.Username) == 0 {
			conf.Username = DefaultConfig.Username
		}
		if len(conf.Password) == 0 {
			conf.Password = DefaultConfig.Password
		}
	}

	return broker
}

func (b Broker) connSub(conf *BrokerConf, messageSubHandler mqtt.MessageHandler) {
	options := mqtt.NewClientOptions().
		AddBroker(conf.Addr).
		SetUsername(conf.Username).
		SetPassword(conf.Password)
	log.Println("Broker Addresses: ", options.Servers)
	options.SetClientID(conf.ClientId())
	options.SetConnectTimeout(time.Duration(conf.ConnectTimeout) * time.Second)
	options.SetAutoReconnect(true)
	options.SetKeepAlive(time.Duration(20) * time.Second)
	options.SetMaxReconnectInterval(time.Duration(5) * time.Second)
	options.SetConnectionLostHandler(func(c mqtt.Client, err_ error) {
		doSub(client, conf.multipleTopics(conf.MultipleTopicQoS), messageSubHandler)
	})
	options.SetOnConnectHandler(func(c mqtt.Client) {
		log.Printf("[client connect state] IsConnected:%v, IsConnectionOpen:%v", c.IsConnected(), c.IsConnectionOpen())
	})

	client = mqtt.NewClient(options)
	doConn(client)
	doSub(client, conf.multipleTopics(conf.MultipleTopicQoS), messageSubHandler)
}

func doConn(client mqtt.Client) {
	for {
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			log.Printf("yomo-source connect error, error: %s \n", token.Error())
			time.Sleep(time.Duration(1) * time.Second)
			continue
		}
		log.Printf("yomo-source connect to broker...")
		break
	}
}

func doSub(client mqtt.Client, topics map[string]byte, messageSubHandler mqtt.MessageHandler) {
	for {
		if token := client.SubscribeMultiple(topics, messageSubHandler); token.Wait() && token.Error() != nil {
			log.Printf("yomo-source SubscribeMultiple error: %s \n", token.Error())
			time.Sleep(500 * time.Millisecond)
			continue
		}
		log.Printf("yomo-souce SubscribeMultiple: %v", topics)
		break
	}
}

func mqttRun(addr string, wait *sync.WaitGroup) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	config, err := broker.ConfigureConfig(os.Args[1:])
	if err != nil {
		log.Fatal("configure broker config error: ", err)
	}

	config.Host, config.Port = getHostPort(addr)

	b, err := broker.NewBroker(config)
	if err != nil {
		log.Fatal("New Broker error: ", err)
	}
	b.Start()
	wait.Done()

	s := waitForSignal()
	log.Println("signal received, broker closed.", s)
}

func waitForSignal() os.Signal {
	signalChan := make(chan os.Signal, 1)
	defer close(signalChan)
	signal.Notify(signalChan, os.Kill, os.Interrupt)
	s := <-signalChan
	signal.Stop(signalChan)
	return s
}
