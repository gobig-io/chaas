package bot

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// Action holds the config, messenger, the script and options to run.
type Action struct {
	config    *Config
	messenger Messenger
	Script    string
	Options   []*Option
}

// UnmarshalText sets the Script to the text string
func (a *Action) UnmarshalText(text []byte) error {
	a.Script = string(text)
	return nil
}

// Init adds the config and messenger to the struct
func (a *Action) Init(config *Config, bot Messenger) {
	a.config = config
	a.messenger = bot
}

// GetOptions calls the `make options` to get the options for the script
func (a *Action) GetOptions(message string) (Options, error) {
	var options Options
	result := a.Make(NewTarget("options", nil, nil))
	if result.Status > 0 {
		return options, fmt.Errorf("%s", result.Error)
	}
	ops := strings.Split(result.Message, "\n")
	for _, op := range ops {
		option, err := NewOption(op, message)
		if err != nil {
			return options, err
		}
		options = append(options, option)
	}
	return options, nil
}

// Make calls the target in the Makefile
func (a *Action) Make(target *Target) *Result {
	cmd := exec.Command("make", target.Call()...)
	cmd.Dir = a.Path()
	cmd.Env = target.Env
	out, err := cmd.CombinedOutput()
	if err != nil {
		return &Result{Status: 1, Message: strings.TrimSpace(string(out)), Error: err.Error()}
	}
	return &Result{Status: 0, Message: strings.TrimSpace(string(out))}
}

// MakeStream calls the target and outputs directly to stdout and stderr
func (a *Action) MakeStream(target *Target) {
	cmd := exec.Command("make", target.Call()...)
	cmd.Dir = a.Path()
	cmd.Env = target.Env
	cmd.Stdout = a.messenger
	cmd.Stderr = a.messenger
	cmd.Run()
}

// Path gets the absolute path for the script
func (a *Action) Path() string {
	p, err := filepath.Abs(filepath.Join(a.config.Actions, a.Script))
	if err != nil {
		fmt.Printf("%s\n", err)
		return ""
	}
	return p
}

// Option holds the key and value along with a Regex for parsing
type Option struct {
	Key   string
	Value string
	Regex *regexp.Regexp
}

// NewOption sets up the Option struct with the key, value and regex
func NewOption(op string, msg string) (*Option, error) {
	re, err := regexp.Compile(`^([^=]+)=(.*)$`)
	if err != nil {
		return &Option{}, err
	}
	p := re.FindStringSubmatch(op)
	if len(p) < 3 {
		return &Option{}, fmt.Errorf("Expected at least 2 in %#v", p)
	}
	key := p[1]
	replacement := p[2]
	replace := fmt.Sprintf("%s ([a-zA-Z0-9-_.!@#*()]+)", replacement)
	re, err = regexp.Compile(replace)
	if err != nil {
		return &Option{}, err
	}
	option := &Option{Key: key, Regex: re}
	if strings.HasSuffix(option.Key, "*") {
		option.Key = strings.TrimSuffix(option.Key, "*")
		replace = fmt.Sprintf("%s (.*)", replacement)
		re, err = regexp.Compile(replace)
		if err != nil {
			return option, err
		}
		option.Regex = re
	}
	rs := re.FindStringSubmatch(msg)
	if len(rs) > 1 {
		option.Value = rs[1]
	}
	return option, nil
}

// Options holds a slice of Option
type Options []*Option

func (o Options) toSlice() []string {
	var ops []string
	for _, op := range o {
		ops = append(ops, op.String())
	}
	return ops
}

func (op *Option) String() string {
	return fmt.Sprintf("%s=\"%s\"", op.Key, op.Value)
}

// Result holds the status, message, error and data from a target
type Result struct {
	Status  int
	Message string
	Error   string
	Data    string
}

// Target holds the Name, Options and Env for a give target
type Target struct {
	Name    string
	Options Options
	Env     []string
}

// NewTarget sets up the Name, Options and Env
func NewTarget(name string, options Options, env []string) *Target {
	return &Target{name, options, env}
}

// Call prepares all the options for a Target
func (t *Target) Call() []string {
	ops := []string{"-s"}
	if t.Name != "" {
		ops = append(ops, t.Name)
	}
	ops = append(ops, t.Options.toSlice()...)
	return ops
}
