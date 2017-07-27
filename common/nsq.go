package common

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nsqio/go-nsq"
)

// CommandHandler listens for commands on a given topic/channel
type CommandHandler struct {
	topic     string
	channel   string
	nsqConfig nsq.Config
	consumer  nsq.Consumer
}

// HandleMessage handles nsq messages on it's topic/channel
func (ch *CommandHandler) HandleMessage(m *nsq.Message) error {
	fmt.Println(m.Body)
	return nil
}

// Start listening for messages on your channel
func (ch *CommandHandler) Start() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-ch.consumer.StopChan:
			return
		case <-sigChan:
			ch.consumer.Stop()
		}
	}
}

func newCommandHandler(topic, channel string) *CommandHandler {
	cfg := nsq.NewConfig()

	consumer, err := nsq.NewConsumer(*topic, *channel, cfg)
	if err != nil {
		log.Fatal(err)
	}

	ch := CommandHandler{
		topic,
		channel,
		cfg,
		consumer,
	}

	consumer.AddHandler(&ch)
	err = consumer.ConnectToNSQLookupd("nsqlookupd")
	if err != nil {
		log.Fatal(err)
	}
	return &ch
}
