package config

import (
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/key"
)

const ConfigDir = "taskyzator"
const configFilename = "config.yaml"

type Key string

func (k Key) prepareToProccess() string {
	var s = strings.ReplaceAll(string(k), "space", " ")
	s = strings.ReplaceAll(s, "↑", "up")
	s = strings.ReplaceAll(s, "↓", "down")
	s = strings.ReplaceAll(s, "←", "left")
	s = strings.ReplaceAll(s, "→", "right")
	return s
}

func (k Key) prepareToDisplay() string {
	var s = strings.ReplaceAll(string(k), " ", "space")
	s = strings.ReplaceAll(s, "up", "↑")
	s = strings.ReplaceAll(s, "down", "↓")
	s = strings.ReplaceAll(s, "left", "←")
	s = strings.ReplaceAll(s, "right", "→")
	return s
}

func (k Key) Key() string {
	return k.prepareToProccess()
}

func (k Key) Binding() key.BindingOpt {
	s := k.prepareToProccess()
	keys := strings.Split(s, ",")
	return key.WithKeys(keys...)
}

func (k Key) Help(help string) key.BindingOpt {
	s := k.prepareToDisplay()
	return key.WithHelp(s, help)
}

func (k Key) Contains(keyName string) bool {
	s := k.prepareToProccess()
	keys := strings.Split(s, ",")
	return slices.Contains(keys, keyName)
}

type Style struct {
	ErrorColor            string `yaml:"error-color"`
	AccentColor           string `yaml:"accent-color"`
	DimmedColor           string `yaml:"dimmed-color"`
	NormalColor           string `yaml:"normal-color"`
	ActiveTaskTextColor   string `yaml:"active-task-text-color"`
	DoneTaskTextColor     string `yaml:"done-task-text-color"`
	ArchivedTaskTextColor string `yaml:"archived-task-text-color"`
}

type Controls struct {
	// Main control
	Quit       Key `yaml:"quit"`
	Apply      Key `yaml:"apply"`
	Cancel     Key `yaml:"cancel"`
	CursorUp   Key `yaml:"cursor-up"`
	CursorDown Key `yaml:"cursor-down"`
	// tasks
	NewTask     Key `yaml:"new-task"`
	DoneTask    Key `yaml:"done-task"`
	ArchiveTask Key `yaml:"archive-task"`
}

type Config struct {
	Controls Controls `yaml:"controls"`
	Style    Style    `yaml:"style"`
}

var defaultConfig = Config{
	Controls: Controls{
		Quit:        "ctrl+q,ctrl+c",
		Apply:       "enter",
		Cancel:      "esc",
		CursorUp:    "up,j",
		CursorDown:  "down,k",
		NewTask:     "ctrl+n",
		DoneTask:    "enter",
		ArchiveTask: "delete",
	},
	Style: Style{
		ErrorColor:            "#F33",
		AccentColor:           "#0cc47e",
		DimmedColor:           "#a3a3a3",
		NormalColor:           "#f2f2f2",
		ActiveTaskTextColor:   "#f7fffc",
		DoneTaskTextColor:     "#6ca690",
		ArchivedTaskTextColor: "#7d7d7d",
	},
}
