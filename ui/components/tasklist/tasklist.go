package tasklist

import (
	"fmt"
	"io"
	"time"

	"github.com/dece2183/taskyzator/config"
	"github.com/dece2183/taskyzator/tasks"
	"github.com/dece2183/taskyzator/ui/model"
	"github.com/dece2183/taskyzator/ui/style"

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

type itemType uint

const (
	_ITEM_TASK itemType = iota
	_ITEM_DONE
	_ITEM_ARCHIVE
)

type Item struct {
	task     *tasks.Task
	itemType itemType
}

func (i Item) FilterValue() string {
	return i.task.Name
}

type ItemDelegate struct{}

func (d ItemDelegate) Height() int {
	return 3
}

func (d ItemDelegate) Spacing() int {
	return 0
}

func (d ItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (d ItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	var name string

	item, ok := listItem.(Item)
	if !ok {
		return
	}

	if item.task != nil {
		var nameStyle lipgloss.Style
		var nameIcon, time string

		name = item.task.Name
		startTime := item.task.BeginTime

		switch item.task.Status {
		case tasks.ACTIVE:
			time = style.DimmedText.Render("start: ") + formatDateTime(startTime)
			nameStyle = activeTaskTextStyle
			nameIcon = " "
		case tasks.DONE:
			time = style.DimmedText.Render("spent: ") + formatDuration(item.task.DoneTime.Sub(startTime)) +
				style.DimmedText.Render(", start: ") + formatDateTime(startTime) +
				style.DimmedText.Render(", done: ") + formatDateTime(item.task.DoneTime)

			nameStyle = doneTaskTextStyle
			nameIcon = style.DoneIcon
		case tasks.ARCHIVED:
			time = style.DimmedText.Render("spent: ") + formatDuration(item.task.DoneTime.Sub(startTime)) +
				style.DimmedText.Render(", start: ") + formatDateTime(startTime) +
				style.DimmedText.Render(", done: ") + formatDateTime(item.task.DoneTime)

			nameStyle = archivedTaskTextStyle
			nameIcon = style.ArchiveIcon
		}

		width := m.Width() - (8 + lipgloss.Width(time))
		nameStyle = nameStyle.Copy().Width(width).MaxWidth(width)
		name = fmt.Sprintf("[%s] %s %s", nameIcon, nameStyle.Render(name), time)
	} else {
		switch item.itemType {
		case _ITEM_DONE:
			name = style.DimmedText.Render("done tasks")
		case _ITEM_ARCHIVE:
			name = style.DimmedText.Render("archive")
		}
	}

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

func New(title string) *Model {
	m := &Model{}

	controls := config.Current.Controls

	m.list = list.New([]list.Item{}, ItemDelegate{}, 512, 512)
	m.list.Title = title
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
	m.sortTasks()

	return m
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) View() string {
	return m.list.View()
}

func (m *Model) Update(message tea.Msg) (*Model, tea.Cmd) {
	var cmd tea.Cmd

	msg, isKeyMessage := message.(tea.KeyMsg)
	if !isKeyMessage {
		return m, nil
	}

	controls := config.Current.Controls
	keypress := msg.String()
	item, isItem := m.list.SelectedItem().(Item)

	switch {
	case controls.NewTask.Contains(keypress):
		cmd = model.Cmd(NEW_TASK)
	case controls.DoneTask.Contains(keypress) && isItem && item.task != nil && item.task.Status == tasks.ACTIVE:
		err := tasks.Done(item.task)
		if err != nil {
			return m, model.Error(err)
		}
		m.sortTasks()
	case controls.UndoneTask.Contains(keypress) && isItem && item.task != nil && item.task.Status == tasks.DONE:
		err := tasks.Undone(item.task)
		if err != nil {
			return m, model.Error(err)
		}
		m.sortTasks()
		m.list.CursorDown()
	case controls.ArchiveTask.Contains(keypress) && isItem:
		if item.task != nil {
			if item.task.Status != tasks.DONE {
				return m, nil
			}

			err := tasks.Archive(item.task)
			if err != nil {
				return m, model.Error(err)
			}
		} else {
			if item.itemType != _ITEM_DONE {
				return m, nil
			}

			items := m.list.Items()
			for i := m.list.Index() + 1; i < len(items)-1; i++ {
				tsk := items[i].(Item)
				if tsk.task == nil {
					break
				}
				tasks.Archive(tsk.task)
			}
		}
		m.sortTasks()
	case controls.UnarchiveTask.Contains(keypress) && isItem && item.task != nil && item.task.Status == tasks.ARCHIVED:
		err := tasks.Unarchive(item.task)
		if err != nil {
			return m, model.Error(err)
		}
		m.sortTasks()
		m.list.CursorDown()
	case controls.DeleteTask.Contains(keypress) && isItem:
		if item.task != nil {
			if item.task.Status != tasks.ARCHIVED {
				return m, nil
			}

			err := tasks.Delete(item.task)
			if err != nil {
				return m, model.Error(err)
			}
		} else {
			if item.itemType != _ITEM_ARCHIVE {
				return m, nil
			}

			err := tasks.DeleteArchived()
			if err != nil {
				return m, model.Error(err)
			}
		}
		m.sortTasks()

	case controls.CursorUp.Contains(keypress):
		m.list, cmd = m.list.Update(msg)
		cmd = tea.Batch(cmd, model.Cmd(CURSOR_UP))
		return m, cmd
	case controls.CursorDown.Contains(keypress):
		m.list, cmd = m.list.Update(msg)
		cmd = tea.Batch(cmd, model.Cmd(CURSOR_DOWN))
		return m, cmd
	default:
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}

	if m.list.Index() >= len(m.list.Items()) || m.list.SelectedItem().(Item).task == nil {
		m.list.CursorUp()
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

func (m *Model) sortTasks() {
	taskList := tasks.List()
	if len(taskList) == 0 {
		m.list.SetItems(nil)
		return
	}

	active_tasks := make([]list.Item, 0, len(taskList)/4)
	done_tasks := make([]list.Item, 0, len(taskList)/4)
	archive_tasks := make([]list.Item, 0, len(taskList)/2)

	for _, item := range taskList {
		switch item.Status {
		case tasks.ACTIVE:
			active_tasks = append(active_tasks, Item{task: item})
		case tasks.DONE:
			done_tasks = append(done_tasks, Item{task: item})
		case tasks.ARCHIVED:
			archive_tasks = append(archive_tasks, Item{task: item})
		}
	}

	itemList := make([]list.Item, 0, len(active_tasks)+len(done_tasks)+len(archive_tasks)+2)
	itemList = append(itemList, active_tasks...)

	itemList = append(itemList, Item{itemType: _ITEM_DONE})
	itemList = append(itemList, done_tasks...)

	itemList = append(itemList, Item{itemType: _ITEM_ARCHIVE})
	itemList = append(itemList, archive_tasks...)

	m.list.SetItems(itemList)
}

func (m *Model) keymap() []key.Binding {
	controls := config.Current.Controls
	bindings := []key.Binding{
		key.NewBinding(controls.NewTask.Binding(), controls.NewTask.Help("new")),
	}

	selectedItem, isItem := m.list.SelectedItem().(Item)
	if len(m.list.Items()) > 0 && isItem {
		if selectedItem.task != nil {
			switch selectedItem.task.Status {
			case tasks.ACTIVE:
				bindings = append(bindings, key.NewBinding(controls.DoneTask.Binding(), controls.DoneTask.Help("done")))
			case tasks.DONE:
				bindings = append(bindings, key.NewBinding(controls.UndoneTask.Binding(), controls.UndoneTask.Help("undone")))
				bindings = append(bindings, key.NewBinding(controls.ArchiveTask.Binding(), controls.ArchiveTask.Help("archive")))
			case tasks.ARCHIVED:
				bindings = append(bindings, key.NewBinding(controls.UnarchiveTask.Binding(), controls.UnarchiveTask.Help("unarchive")))
				bindings = append(bindings, key.NewBinding(controls.DeleteTask.Binding(), controls.DeleteTask.Help("permanet delete")))
			}
		} else {
			switch selectedItem.itemType {
			case _ITEM_DONE:
				bindings = append(bindings, key.NewBinding(controls.ArchiveTask.Binding(), controls.ArchiveTask.Help("archive all")))
			case _ITEM_ARCHIVE:
				bindings = append(bindings, key.NewBinding(controls.DeleteTask.Binding(), controls.DeleteTask.Help("permanet delete all")))
			}
		}
	}

	return bindings
}

func formatDuration(d time.Duration) string {
	return style.NormalText.Render(fmt.Sprintf(
		"%02d:%02d",
		int(d.Hours()), int(d.Minutes())-int(d.Hours())*60,
	))
}

func formatDateTime(t time.Time) string {
	return style.NormalText.Render(fmt.Sprintf(
		"%02d.%02d.%04d - %02d:%02d",
		t.Day(), t.Month(), t.Year(),
		t.Hour(), t.Minute(),
	))
}
