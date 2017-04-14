package terminal

import (
	"bufio"
	"fmt"
	"os"

	"github.com/gobig-io/chaas"
)

func New(c *bot.Config) *Terminal {
	return &Terminal{c}
}

type Terminal struct {
	config *bot.Config
}

func (t *Terminal) Read(data []byte) (int, error) {
	fmt.Printf("%s> ", t.config.Name)
	n, err := os.Stdin.Read(data)
	return n, err
}

func (t *Terminal) Listen(msg *bot.Message) error {
	for {
		scanner := bufio.NewScanner(t)
		err := msg.Scan(scanner)
		return err
	}
}

func (t *Terminal) Write(data []byte) (int, error) {
	n, err := os.Stdout.Write(data)
	return n, err
}

func (t *Terminal) Respond(results []*bot.Result) error {
	writer := bufio.NewWriter(t)
	for _, r := range results {
		if r.Status > 0 {
			if _, err := writer.WriteString(r.Error + "\n"); err != nil {
				return err
			}
			continue
		}
		if _, err := writer.WriteString(r.Message + "\n"); err != nil {
			return err
		}
	}
	writer.Flush()
	return nil
}
