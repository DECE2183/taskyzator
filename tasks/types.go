package tasks

import "time"

const tasksFilename = "tasklist.yaml"

type Status int

const (
	ACTIVE Status = iota
	DONE
	ARCHIVED
)

type Task struct {
	Name      string    `yaml:"name"`
	Status    Status    `yaml:"status"`
	BeginTime time.Time `yaml:"begin-time"`
	DoneTime  time.Time `yaml:"done-time"`
}
