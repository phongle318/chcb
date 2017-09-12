package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	fptai "github.com/fpt-corp/fptai-sdk-go"
	"github.com/michlabs/fbbot"
)

const (
	FPTAI_TOKEN = "DCYgINEOZCbxtsgAEO4vFnpArEsZqLse"

	FB_PAGE_ACCESS_TOKEN = "EAALFPaj8iGIBAEHul76ktvUunolgiSuP9pYEZCHqmo8vQOFZChgvmSTrURIpRU6hzoeIeRUM1b6kknfZAit36OYQF8dlOSFDUcPqZAR9eaqnZCYdq0BLuVMowPlo9ppwb9UDd21GyLN8UbzdGh6uqho8bOk3e0mctyVEEAZC1uuwZDZD"
	FB_VERIFY_TOKEN      = "PhongLH318"
)

var client *fptai.Client
var PORT int = 1203
var ErrMsg string = "Sory, But I don't understand what are you saying."

type BjjNerd struct {
	Answers []Answer
}

type Answer struct {
	Intent string `json:"intent"`
	Text   string `json:"answer"`
}

func (t *BjjNerd) Load(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &t.Answers); err != nil {
		return err
	}
	return nil
}

func (t *BjjNerd) Answer(intent string) string {
	for _, answer := range t.Answers {
		if answer.Intent == intent {
			return answer.Text
		}
	}
	return ErrMsg
}

func (t *BjjNerd) HandleMessage(bot *fbbot.Bot, msg *fbbot.Message) {
	bot.TypingOn(msg.Sender)
	resp, err := client.RecognizeIntents(msg.Text)
	if err != nil || len(resp.Intents) == 0 {
		bot.SendText(msg.Sender, ErrMsg)
		return
	}
	intent := resp.Intents[0].Name
	bot.SendText(msg.Sender, t.Answer(intent))
}

func init() {
	client = fptai.NewClient(FPTAI_TOKEN)
}

func main() {
	var bjjNerd BjjNerd
	if err := bjjNerd.Load("answers.json"); err != nil {
		log.Fatal(err)
	}

	bot := fbbot.New(PORT, FB_VERIFY_TOKEN, FB_PAGE_ACCESS_TOKEN)
	bot.AddMessageHandler(&bjjNerd)
	bot.Run()
}
