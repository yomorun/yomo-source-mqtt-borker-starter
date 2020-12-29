package starter

import (
	"fmt"
	"time"
)

type BrokerConf struct {
	Addr           string
	Topics         []string
	Username       string
	Password       string
	ConnectTimeout int
	FailureTimes   int
}

func (c BrokerConf) ClientId() string {
	return fmt.Sprintf("yomo-source-sub-%d", time.Now().Unix())
}

func (c BrokerConf) multipleTopics() map[string]byte {
	topics := make(map[string]byte)
	for _, topic := range c.Topics {
		topics[topic] = byte(1)
	}
	return topics
}

var DefaultConfig = &BrokerConf{
	Addr:           "localhost:1883",
	Topics:         []string{"NOISE"},
	Username:       "admin",
	Password:       "public",
	ConnectTimeout: 0,
	FailureTimes:   5,
}
