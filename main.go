package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token string
)

type Responses struct {
	Responses []string `json:"responses"`
}

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	//Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	hc, err := runHateDetectionCmd(m.Content)
	for _, i := range hc.Classes {
		if i.ClassName == "hate_speech" {
			if i.Confidence > 0.15 {
				res, _ := ioutil.ReadFile("responses.json")

				var responses Responses
				err := json.Unmarshal(res, &responses)
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, "Hate Speech detected!")
				}
				rand.Seed(time.Now().Unix()) // initialize global pseudo random generator
				message := responses.Responses[rand.Intn(len(responses.Responses))]
				s.ChannelMessageSend(m.ChannelID, message)
			}
		}
	}
	if err != nil {
		fmt.Println(err)
		return
	}
}
