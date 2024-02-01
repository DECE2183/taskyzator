package mainpage

import (
	"github.com/dece2183/taskyzator/config"
	"github.com/dece2183/taskyzator/ui/components/input"
	"github.com/dece2183/taskyzator/ui/components/tasklist"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.design/x/clipboard"
)

type Model struct {
	program       *tea.Program
	width, height int
	inNewTaskMode bool
	tasklist      *tasklist.Model
	taskNameInput *input.Model
}

func New() *Model {
	m := &Model{}
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	m.program = p
	m.tasklist = tasklist.New("taskyzator")
	m.taskNameInput = input.New("enter new task name:")
	return m
}

//
// model.Model interface implementation
//

func (m *Model) Run() error {
	err := clipboard.Init()
	if err != nil {
		return err
	}

	_, err = m.program.Run()
	return err
}

func (m *Model) Send(msg tea.Msg) {
	go m.program.Send(msg)
}

//
// tea.Model interface implementation
//

func (m *Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := message.(type) {
	case tea.WindowSizeMsg:
		m.resize(msg.Width, msg.Height)
		return m, tea.ClearScreen

	case tea.KeyMsg:
		controls := config.Current.Controls
		keypress := msg.String()

		if controls.Quit.Contains(keypress) {
			return m, tea.Quit
		}

		if m.inNewTaskMode {
			m.taskNameInput, cmd = m.taskNameInput.Update(msg)
		} else {
			m.tasklist, cmd = m.tasklist.Update(msg)
		}

		cmds = append(cmds, cmd)

	case tasklist.TasklistControl:
		if msg == tasklist.NEW_TASK {
			m.inNewTaskMode = true
			cmd = m.taskNameInput.Focus()
			cmds = append(cmds, cmd)
		}

	case input.InputControl:
		str := msg.String()
		if len(str) > 0 {
			cmd = m.tasklist.NewTask(str)
			cmds = append(cmds, cmd)
		}
		m.inNewTaskMode = false

	default:
		if m.inNewTaskMode {
			m.taskNameInput, cmd = m.taskNameInput.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	if m.inNewTaskMode {
		return m.viewNewTaskInput()
	}
	return m.viewTaskList()
}

//
// private methods
//

func (m *Model) viewNewTaskInput() string {
	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		m.taskNameInput.View(),
	)
}

func (m *Model) viewTaskList() string {
	return m.tasklist.View()
}

func (m *Model) resize(width, height int) {
	m.width, m.height = width, height
	m.tasklist.SetSize(width-1, height-1)
}
