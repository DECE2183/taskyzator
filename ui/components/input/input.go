package input

import (
	"taskyzator/config"
	"taskyzator/ui/model"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type InputControl string

func (i InputControl) String() string {
	return string(i)
}

type helpKeyMap struct {
	Cancel key.Binding
	Apply  key.Binding
}

func (k helpKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Cancel, k.Apply}
}

func (k helpKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Cancel, k.Apply},
	}
}

var helpMap = helpKeyMap{
	Cancel: key.NewBinding(
		config.Current.Controls.Cancel.Binding(),
		config.Current.Controls.Cancel.Help("cancel"),
	),
	Apply: key.NewBinding(
		config.Current.Controls.Apply.Binding(),
		config.Current.Controls.Apply.Help("ok"),
	),
}

type Model struct {
	input textinput.Model
	help  help.Model
	title string
}

func New(title string) *Model {
	m := &Model{
		title: title,
	}

	m.help = help.New()

	m.input = textinput.New()
	m.input.Width = 80
	m.input.CharLimit = 256

	return m
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) View() string {
	body := lipgloss.JoinVertical(
		lipgloss.Left,
		inputTitleStyle.Render(m.title),
		inputTextStyle.Render(m.input.View()),
		m.help.View(helpMap),
	)
	body = inputBoxStyle.Render(body)
	return body
}

func (m *Model) Update(message tea.Msg) (*Model, tea.Cmd) {
	var cmd tea.Cmd

	msg, isKeyMessage := message.(tea.KeyMsg)
	if !isKeyMessage {
		return m, nil
	}

	controls := config.Current.Controls
	keypress := msg.String()

	switch {
	case controls.Apply.Contains(keypress):
		str := m.input.Value()
		cmd = model.Cmd(InputControl(str))
		m.input.Reset()
		m.input.Blur()
	case controls.Cancel.Contains(keypress):
		cmd = model.Cmd(InputControl(""))
		m.input.Reset()
		m.input.Blur()
	default:
		m.input, cmd = m.input.Update(message)
	}

	return m, cmd
}

func (m *Model) Focus() tea.Cmd {
	return m.input.Focus()
}
