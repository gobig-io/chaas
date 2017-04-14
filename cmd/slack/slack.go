package main

import (
	"log"
	"os"

	"github.com/gobig-io/chaas"
	"github.com/gobig-io/chaas/messengers/slack"
)

func main() {
	f := NewFlags(os.Args)
	config := bot.NewConfig(f.Name, f.Directions)
	config.ID = f.ID
	config.Token = f.Token
	b := slack.New(config)
	directions, err := bot.NewDirections(config)
	if err != nil {
		log.Fatal(err)
	}
	for {
		msg := bot.NewMessage(config)
		if err := b.Listen(msg); err != nil {
			//log.Println(err)
			continue
		}
		go func(msg *bot.Message) {
			if err := msg.Process(b, directions); err != nil {
				log.Println(err)
			}
		}(msg)
	}
}
