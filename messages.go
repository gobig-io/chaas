package bot

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"unicode/utf8"
)

// Message holds the config, words, directions, User, Email, and Time
type Message struct {
	config     *Config
	words      []string
	directions Directions
	User       string
	Email      string
	Time       time.Time
}

// NewMesssage sets up a new message with a given config
func NewMessage(c *Config) *Message {
	return &Message{config: c}
}

// AddProfile adds the user, email, time to the Message
func (m *Message) AddProfile(user, email string, date time.Time) {
	m.User = user
	m.Email = email
	m.Time = date
}

func (m *Message) String() string {
	return strings.Join(m.words, " ")
}

// AddWord adds a new parsed word to the message
func (m *Message) AddWord(word string) {
	m.words = append(m.words, word)
}

// Scan goes over the entire message to add words
func (m *Message) Scan(scanner *bufio.Scanner) error {
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		m.AddWord(scanner.Text())
	}
	// fmt.Println("Scan")
	// fmt.Printf("words: %#v\n", m.words)
	if err := scanner.Err(); err != nil && err.Error() != "atEOD" {
		return err
	}
	return nil
}

// ScanWords is used to split the text into words
func (m *Message) ScanWords(data []byte, atEOF bool) (advance int, token []byte, err error) {
	start := 0
	for width := 0; start < len(data); start += width {
		var r rune
		r, width = utf8.DecodeRune(data[start:])
		if !isSpace(r) {
			fmt.Println("no space")
			break
		}
	}
	for width, i := 0, start; i < len(data); i += width {
		var r rune
		r, width = utf8.DecodeRune(data[i:])
		if isSpace(r) {
			fmt.Println("found a space")
			return i + width, data[start:i], nil
		}
	}
	if atEOF && len(data) == 0 {
		return 0, nil, fmt.Errorf("%s", "empty")
	}
	if atEOF && len(data) > start {
		fmt.Println("atEOF")
		return len(data), data[start:], nil
	}
	fmt.Println("want more data")
	return len(data), data[start:], nil
}

// Process initializes the Action, adds the user, email, time to the env
// gets the options, calls the intro target, and finally streams the target
// results to the messenger
func (m *Message) Process(bot Messenger, directions Directions) error {
	directions, err := m.processDirections(directions)
	if err != nil {
		return err
	}
	env := append(os.Environ(), []string{
		"user=" + m.User,
		"email=" + m.Email,
		"time=" + m.Time.Format(time.RFC3339),
	}...)
	for _, direction := range directions {
		for _, a := range direction.Actions {
			var ops Options
			a.Init(m.config, bot)
			if ops, err = a.GetOptions(m.String()); err != nil {
				log.Println("Error: make options:", err.Error())
			}
			env = append(env, "target="+direction.Target)
			result := a.Make(NewTarget("intro", ops, env))
			if result.Status == 0 {
				if err = bot.Respond([]*Result{result}); err != nil {
					log.Println("Error: make intro:", err.Error())
				}
			}
			a.MakeStream(NewTarget(direction.Target, ops, env))
		}
	}
	return nil
}

func (m *Message) processDirections(directions Directions) (Directions, error) {
	if len(m.words) < 1 {
		return nil, fmt.Errorf("%s", "No words found")
	}
	for _, direction := range directions {
		for _, text := range direction.Words {
			var word = m.words[0]
			if word == m.config.ID {
				word = m.words[1]
			}
			if strings.Title(text) == strings.Title(word) {
				direction.Target = text
				m.directions = append(m.directions, direction)
			}
		}
	}
	return m.directions, nil
}
