package tasks

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/dece2183/taskyzator/config"

	"gopkg.in/yaml.v3"
)

type TaskList []*Task

var taskList TaskList

func init() {
	var err error
	taskList, err = load()
	if err != nil {
		save(taskList)
	}
}

func getDir() (string, error) {
	userDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(userDir, ".config", config.ConfigDir)
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		return "", err
	}

	return configDir, nil
}

func load() (TaskList, error) {
	configDir, err := getDir()
	if err != nil {
		return TaskList{}, err
	}

	configContent, err := os.ReadFile(filepath.Join(configDir, tasksFilename))
	if err != nil {
		return TaskList{}, err
	}

	var newList []Task
	err = yaml.Unmarshal(configContent, &newList)
	if err != nil {
		return TaskList{}, err
	}

	tasks := make(TaskList, len(newList))
	for i := range newList {
		tasks[i] = &newList[i]
	}

	return tasks, nil
}

func save(list TaskList) error {
	configDir, err := getDir()
	if err != nil {
		return err
	}

	file, err := os.OpenFile(filepath.Join(configDir, tasksFilename), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := yaml.NewEncoder(file)
	enc.SetIndent(4)
	err = enc.Encode(list)
	if err != nil {
		return err
	}

	return nil
}

func List() TaskList {
	return taskList
}

func Save() error {
	return save(taskList)
}

func Append(name string) (*Task, error) {
	newTask := Task{
		Name:      name,
		Status:    ACTIVE,
		BeginTime: time.Now(),
	}

	taskList = append(TaskList{&newTask}, taskList...)
	return &newTask, Save()
}

func Done(task *Task) error {
	task.Status = DONE
	task.DoneTime = time.Now()
	return Save()
}

func Undone(task *Task) error {
	task.Status = ACTIVE
	return Save()
}

func Archive(task *Task) error {
	task.Status = ARCHIVED
	return Save()
}

func Unarchive(task *Task) error {
	task.Status = DONE
	return Save()
}

func Delete(task *Task) error {
	for i := range taskList {
		if taskList[i] == task {
			switch i {
			case 0:
				taskList = taskList[i+1:]
			case len(taskList) - 1:
				taskList = taskList[:len(taskList)-1]
			default:
				taskList = append(taskList[:i], taskList[i+1:]...)
			}
			return Save()
		}
	}
	return fmt.Errorf("task not found")
}

func DeleteArchived() error {
	i := 0
	for i < len(taskList) {
		if taskList[i].Status != ARCHIVED {
			continue
		}

		if len(taskList) == 1 {
			taskList = taskList[0:0]
			break
		}

		switch i {
		case 0:
			taskList = taskList[i+1:]
		case len(taskList) - 1:
			taskList = taskList[:len(taskList)-1]
		default:
			taskList = append(taskList[:i], taskList[i+1:]...)
		}
		i++
	}
	return Save()
}
