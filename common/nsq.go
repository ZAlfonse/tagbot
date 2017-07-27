package common

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nsqio/go-nsq"
)

// CommandHandler listens for commands on a given topic/channel
type CommandHandler struct {
	topic   string
	channel string
	handler nsq.Handler
}

// Start listening for messages on your channel
func (ch *CommandHandler) Start() {
	cfg := nsq.NewConfig()

	consumer, err := nsq.NewConsumer(ch.topic, ch.channel, cfg)
	if err != nil {
		log.Fatal(err)
	}

	consumer.AddHandler(ch.handler)

	err = consumer.ConnectToNSQLookupd("nsqlookupd:4161")
	if err != nil {
		log.Fatal(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-consumer.StopChan:
			return
		case <-sigChan:
			consumer.Stop()
		}
	}
}

//NewCommandHandler ...
func NewCommandHandler(topic, channel string, handler nsq.Handler) *CommandHandler {
	return &CommandHandler{
		topic,
		channel,
		handler,
	}
}
