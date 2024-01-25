package tasklist

import (
	"fmt"
	"io"
	"taskyzator/config"
	"taskyzator/tasks"
	"taskyzator/ui/model"
	"taskyzator/ui/style"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TasklistControl uint

const (
	CURSOR_UP TasklistControl = iota
	CURSOR_DOWN
	NEW_TASK
)

type Item struct {
	task *tasks.Task
}

func (i Item) FilterValue() string {
	return i.task.Name
}

type ItemDelegate struct{}

func (d ItemDelegate) Height() int {
	return 4
}

func (d ItemDelegate) Spacing() int {
	return 0
}

func (d ItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (d ItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(Item)
	if !ok {
		return
	}

	var nameStyle lipgloss.Style
	var nameIcon string
	name := item.task.Name

	switch item.task.Status {
	case tasks.ACTIVE:
		nameStyle = activeTaskTextStyle
		nameIcon = " "
	case tasks.DONE:
		nameStyle = doneTaskTextStyle
		nameIcon = style.DoneIcon
	case tasks.ARCHIVED:
		nameStyle = archivedTaskTextStyle
		nameIcon = style.ArchiveIcon
	}

	startTime := item.task.BeginTime

	name = fmt.Sprintf(
		"[%s] %s %02d:%02d:%02d",
		nameIcon, nameStyle.Copy().Width(m.Width()-16).Render(name),
		startTime.Hour(), startTime.Minute(), startTime.Second(),
	)

	if index == m.Index() {
		fmt.Fprint(w, selectedItemStyle.Render(name))
	} else {
		fmt.Fprint(w, itemStyle.Render(name))
	}
}

type Model struct {
	list          list.Model
	width, height int
}

func New() Model {
	m := Model{}

	controls := config.Current.Controls

	m.list = list.New([]list.Item{}, ItemDelegate{}, 512, 512)
	m.list.Title = "tasks"
	m.list.Styles.Title = titleStyle
	m.list.KeyMap = list.KeyMap{
		CursorUp:   key.NewBinding(controls.CursorUp.Binding(), controls.CursorUp.Help("up")),
		CursorDown: key.NewBinding(controls.CursorDown.Binding(), controls.CursorDown.Help("down")),
	}
	m.list.AdditionalShortHelpKeys = m.keymap

	items := make([]list.Item, len(tasks.List()))
	for i, task := range tasks.List() {
		items[i] = Item{task: task}
	}
	m.list.SetItems(items)

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	return m.list.View()
}

func (m Model) Update(message tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	msg, isKeyMessage := message.(tea.KeyMsg)
	if !isKeyMessage {
		return m, nil
	}

	controls := config.Current.Controls
	keypress := msg.String()
	var task *tasks.Task
	if len(m.list.Items()) > 0 {
		task = m.list.SelectedItem().(Item).task
	}

	switch {
	case controls.NewTask.Contains(keypress):
		cmd = model.Cmd(NEW_TASK)
	case controls.DoneTask.Contains(keypress) && task != nil && task.Status == tasks.ACTIVE:
		err := tasks.Done(task)
		if err != nil {
			return m, model.Error(err)
		}
	case controls.ArchiveTask.Contains(keypress) && task != nil && task.Status == tasks.DONE:
		err := tasks.Archive(task)
		if err != nil {
			return m, model.Error(err)
		}
	case controls.CursorUp.Contains(keypress):
		m.list, cmd = m.list.Update(msg)
		cmd = tea.Batch(cmd, model.Cmd(CURSOR_UP))
	case controls.CursorDown.Contains(keypress):
		m.list, cmd = m.list.Update(msg)
		cmd = tea.Batch(cmd, model.Cmd(CURSOR_DOWN))
	default:
		m.list, cmd = m.list.Update(msg)
	}

	return m, cmd
}

func (m *Model) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.list.SetSize(m.width, m.height)
}

func (m *Model) NewTask(name string) tea.Cmd {
	newTask, err := tasks.Append(name)
	if err != nil {
		return model.Error(err)
	}
	return m.list.InsertItem(0, Item{task: newTask})
}

func (m *Model) keymap() []key.Binding {
	controls := config.Current.Controls
	bindings := []key.Binding{
		key.NewBinding(controls.NewTask.Binding(), controls.NewTask.Help("new")),
	}

	if len(m.list.Items()) > 0 {
		selectedItem := m.list.SelectedItem().(Item)
		switch selectedItem.task.Status {
		case tasks.ACTIVE:
			bindings = append(bindings, key.NewBinding(controls.DoneTask.Binding(), controls.DoneTask.Help("done")))
		case tasks.DONE:
			bindings = append(bindings, key.NewBinding(controls.ArchiveTask.Binding(), controls.ArchiveTask.Help("archive")))
		}
	}

	return bindings
}
