package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

// Directions holds a slice of Direction
type Directions []*Direction

func (ds Directions) String() string {
	var output string
	for _, d := range ds {
		output += fmt.Sprint("texts:")
		output += strings.Join(d.Words, ",")
		output += fmt.Sprint("actions:")
		for _, a := range d.Actions {
			output += fmt.Sprintf("%#v\n", a)
		}
	}
	return output
}

// Direction holds a target, words, and actions
type Direction struct {
	Target  string
	Words   []string
	Actions []*Action
}

// NewDirections reads the config and sets up the Directions
func NewDirections(config *Config) (Directions, error) {
	data, err := ioutil.ReadFile(config.File())
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, config); err != nil {
		return nil, err
	}
	return config.Directions, nil
}
