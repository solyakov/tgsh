package main

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	flags "github.com/jessevdk/go-flags"
)

func run(request *tgbotapi.Message, shell string, replies chan<- tgbotapi.MessageConfig) {
	cmd := exec.Command(shell)
	cmd.Stdin = strings.NewReader(request.Text)
	res, err := cmd.CombinedOutput()

	if err != nil {
		res = append(res, []byte("\nError: "+err.Error())...)
	}

	if len(res) == 0 {
		res = []byte("No output")
	}

	if len(res) > 4096 {
		t := []byte("\nOutput truncated")
		res = append(res[:4096-len(t)], t...)
	}

	msg := tgbotapi.NewMessage(request.Chat.ID, string(res))
	msg.ReplyToMessageID = request.MessageID

	replies <- msg
}

func main() {
	var opts struct {
		Token string `long:"token" short:"t" env:"TGSH_TOKEN" description:"Telegram Bot Token" required:"true"`
		Shell string `long:"shell" short:"s" env:"TGSH_SHELL" description:"Shell to execute commands" default:"/bin/bash"`
		User  int64  `long:"user" short:"u" env:"TGSH_USER" description:"Telegram User ID allowed to execute commands" required:"true"`
		Debug bool   `long:"debug" short:"d" env:"TGSH_DEBUG" description:"Enable debug mode"`
	}
	if _, err := flags.Parse(&opts); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			return
		}
		log.Fatalf("Error parsing flags: %s", err)
	}

	bot, err := tgbotapi.NewBotAPI(opts.Token)
	if err != nil {
		log.Fatalf("Error initializing bot: %s", err)
	}

	bot.Debug = opts.Debug

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	u.AllowedUpdates = []string{tgbotapi.UpdateTypeMessage}
	updates := bot.GetUpdatesChan(u)
	replies := make(chan tgbotapi.MessageConfig, 10)

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("Ready, authorized for user ID %d using shell %s", opts.User, opts.Shell)

	for {
		select {
		case update := <-updates:
			if update.Message == nil {
				continue
			}
			if update.Message.From.ID != opts.User {
				log.Printf("Unauthorized user: %s (%d)", update.Message.From.UserName, update.Message.From.ID)
				continue
			}
			log.Printf("Received (%d): %s", update.Message.MessageID, update.Message.Text)
			go run(update.Message, opts.Shell, replies)
		case msg := <-replies:
			if _, err := bot.Send(msg); err != nil {
				log.Fatalf("Failed to send message: %s", err)
			}
			log.Printf("Sent (%d): %s", msg.ReplyToMessageID, msg.Text)
		case <-exit:
			log.Println("Shutting down")
			return
		}
	}
}
