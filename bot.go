package main

import (
    "os"
    "os/signal"
    "syscall"
    "strings"
    "net/http"
    "net/url"
    "regexp"
    "io/ioutil"
    "encoding/json"
    "bytes"

    "github.com/zalfonse/tagbot/common"
    "github.com/zalfonse/lumber"
    "github.com/bwmarrin/discordgo"
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
  command_regex, _ := regexp.Compile("[^a-zA-Z0-9]+")
  args_regex, _ := regexp.Compile("[^a-zA-Z0-9]+")

  split := strings.SplitN(line, " ", 2)
  cmd := command_regex.ReplaceAllString(split[0][1:], "")
  args := ""
  if len(split) > 1 {
    args = args_regex.ReplaceAllString(split[1], "")
  }

  return common.Command{cmd, args}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
  if m.Author.ID == s.State.User.ID {
    return
  }
  if m.Content == "!roll" {
    return
  }

  if strings.HasPrefix(m.Content, "!") {
    command := safeCommand(m.Content)

    logger.Info("Command recieved: [" + command.Name + "] with args [" + command.Args +"]")

    command_body, _ := json.Marshal(command)
    resp, err := http.Post("http://" + command.Name + "/execute?args=" + url.QueryEscape(command.Args), "application/json", bytes.NewBuffer(command_body))
    if err != nil {
      s.ChannelMessageSend(m.ChannelID, "Unknown command: " + command.Name)
      logger.Error("Error: ["+ err.Error() + "]")
      return
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    logger.Info("Executed: " + command.Name + ". Got: " + string(body))
    s.ChannelMessageSend(m.ChannelID, string(body))
  }
}
