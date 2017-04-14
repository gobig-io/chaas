package bot

import (
	"os"
	"path/filepath"
)

// Config holds the ID, Name, Directions, Actions, Token, and file for the Bot
type Config struct {
	ID         string
	Name       string
	Directions Directions
	Actions    string
	Token      string
	file       string
}

// NewConfig sets up the name and file
func NewConfig(name string, file string) *Config {
	return &Config{Name: name, file: file}
}

// File returns the full filepath for the config
func (c *Config) File() string {
	file := c.file
	if file == "" {
		return file
	}
	if string(file[0]) != "/" {
		wd, err := os.Getwd()
		if err != nil {
			return ""
		}
		file = filepath.Join(wd, file)
	}
	return file
}
