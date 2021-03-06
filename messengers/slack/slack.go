package slack

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gobig-io/chaas"

	"golang.org/x/net/websocket"
)

var counter uint64

// New contructs a Slack Messenger
func New(c *bot.Config) *Slack {
	slack := &Slack{config: c}
	if err := slack.connect(); err != nil {
		log.Println(err)
	}
	return slack
}

//Slack is a Messenger
type Slack struct {
	ID      string
	config  *bot.Config
	conn    *websocket.Conn
	message *Message
	users   users
}

type users struct {
	Members []struct {
		ID      string
		Name    string
		Profile struct {
			Email string
		}
	} `json:"members"`
}

// Message is a slack message
type Message struct {
	Type    string `json:"type"`
	Channel string `json:"channel"`
	User    string `json:"user"`
	Time    tm     `json:"ts"`
	Text    string `json:"text"`
}

type rtm struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
	URL   string `json:"url"`
	Self  struct {
		ID string `json:"id"`
	} `json:"self"`
}

type tm string

func (t tm) Time() time.Time {
	i, err := strconv.ParseInt(string(t), 10, 64)
	if err != nil {
		return time.Now()
	}
	return time.Unix(i, 0)
}

func (m *Slack) connect() error {
	url := fmt.Sprintf("https://slack.com/api/users.list?token=%s", m.config.Token)
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("API request failed with code %d", res.StatusCode)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return err
	}
	var us users
	err = json.Unmarshal(body, &us)
	if err != nil {
		return err
	}
	url = fmt.Sprintf("https://slack.com/api/rtm.start?token=%s", m.config.Token)
	res, err = http.Get(url)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("API request failed with code %d", res.StatusCode)
	}
	body, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return err
	}
	var message rtm
	err = json.Unmarshal(body, &message)
	if err != nil {
		return err
	}
	if !message.Ok {
		return fmt.Errorf("Slack error: %s", message.Error)
	}
	m.conn, err = websocket.Dial(message.URL, "", "https://api.slack.com/")
	if err != nil {
		return err
	}
	m.ID = message.Self.ID
	m.users = us
	return nil
}

func (m *Slack) Messenger() bot.Messenger {
	return &SlackMessenger{
		b:    m,
		conn: m.conn,
	}
}

type SlackMessenger struct {
	b       *Slack
	channel string
	conn    *websocket.Conn
	message *Message
}

//UserAndEmail returns the user,email
func (m *SlackMessenger) UserAndEmail() (string, string) {
	for _, user := range m.b.users.Members {
		if m.message.User == user.ID {
			return user.Name, user.Profile.Email
		}
	}
	return "", ""
}

func (m *SlackMessenger) Read(b []byte) (int, error) {
	return 0, nil
}

func (m *SlackMessenger) read() (*strings.Reader, error) {
	var msg Message
	if err := websocket.JSON.Receive(m.conn, &msg); err != nil {
		return nil, err
	}
	if msg.Type == "message" {
		m.message = &msg
		return strings.NewReader(msg.Text), nil
	}
	return nil, fmt.Errorf("Unhandled Message Type: %+v", msg)
}

// Listen reads and scans message
func (m *SlackMessenger) Listen(msg *bot.Message) error {
	for {
		reader, err := m.read()
		if err != nil {
			return err
		}
		scanner := bufio.NewScanner(reader)
		user, email := m.UserAndEmail()
		msg.AddProfile(user, email, m.message.Time.Time())
		return msg.Scan(scanner)
	}
}

// Write sends the messange
func (m *SlackMessenger) Write(data []byte) (int, error) {
	time.Sleep(500 * time.Millisecond)
	m.message.User = m.b.config.ID
	m.message.Text = string(data)
	m.message.Time = tm(strconv.FormatInt(time.Now().UTC().Unix(), 10))
	log.Printf("message: %#v\n", m.message)
	if err := websocket.JSON.Send(m.conn, m.message); err != nil {
		fmt.Println(err)
	}
	return len(data), nil
}

// Respond sends the bot results to the writer
func (m *SlackMessenger) Respond(results []*bot.Result) error {
	writer := bufio.NewWriter(m)
	for _, r := range results {
		if r.Status > 0 {
			if _, err := writer.WriteString(r.Message + "\n"); err != nil {
				return err
			}
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
