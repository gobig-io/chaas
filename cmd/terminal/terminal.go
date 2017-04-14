package main

import (
	"log"

	"github.com/gobig-io/chaas"
	"github.com/gobig-io/chaas/messengers/terminal"
)

func main() {
	config := bot.NewConfig("chaas", "directions.json")
	terminal := terminal.New(config)
	for {
		msg := bot.NewMessage(config)
		if err := terminal.Listen(msg); err != nil {
			log.Println(err)
			continue
		}
		responses, err := msg.Process()
		if err != nil {
			log.Println(err)
			continue
		}
		if err := terminal.Respond(responses); err != nil {
			log.Println(err)
		}
	}
}
