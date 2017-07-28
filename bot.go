package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/zalfonse/lumber"
	"github.com/zalfonse/tagbot/common"
)

var logger *lumber.Logger

func main() {
	logger = lumber.NewLogger(lumber.TRACE)
	botToken, exists := os.LookupEnv("BOT_TOKEN")
	if !exists {
		logger.Error("No Bot Token set (Expected as environment variable BOT_TOKEN). Exiting.")
		return
	}

	discord, err := discordgo.New("Bot " + botToken)
	if err != nil {
		logger.Error("Error creating Discord session: ", err)
		return
	}

	discord.AddHandler(messageCreate)

	err = discord.Open()
	if err != nil {
		logger.Error("Error opening connection to discord: ", err)
		return
	}

	logger.Info("TagBot is now running! Ctrl + c to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	discord.Close()
}

func safeCommand(line string) common.Command {
	commandRegex, _ := regexp.Compile("[^a-zA-Z0-9]+")
	argsRegex, _ := regexp.Compile("[^a-zA-Z0-9\\s]+")

	split := strings.SplitN(line, " ", 2)
	cmd := commandRegex.ReplaceAllString(split[0][1:], "")
	args := ""
	if len(split) > 1 {
		args = argsRegex.ReplaceAllString(split[1], "")
	}

	return common.Command{
		Name: cmd,
		Args: args,
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "&") {
		command := safeCommand(m.Content)

		logger.Info("Command recieved: [" + command.Name + "] with args [" + command.Args + "]")

		commandBody, _ := json.Marshal(command)
		resp, err := http.Post("http://"+command.Name+"/execute", "application/json", bytes.NewBuffer(commandBody))
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Unknown command: "+command.Name)
			logger.Error("Error: [" + err.Error() + "]")
			return
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)

		var response common.Response
		if err := json.Unmarshal(body, &response); err != nil {
			logger.Error("Error: [" + err.Error() + "]")
			return
		}
		logger.Info("Executed: " + response.Command.Name + ". Got: " + string(body))
		for _, answer := range response.Answers {
			s.ChannelMessageSend(m.ChannelID, answer)
		}
	}
}
