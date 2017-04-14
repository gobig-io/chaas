package main

import (
	"flag"
	"os"
)

func NewFlags(args []string) *Flags {
	f := &Flags{}
	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	fs.StringVar(&f.Name, "name", "chaas", "Slack Bot Name")
	fs.StringVar(&f.ID, "id", "", "Slack Bot User ID")
	fs.StringVar(&f.Token, "token", "", "Slack Bot API Key")
	fs.StringVar(&f.Directions, "directions", "directions.json", "Path to directions.json")
	fs.Parse(os.Args[1:])
	return f
}

type Flags struct {
	ID         string
	Token      string
	Name       string
	Directions string
}
