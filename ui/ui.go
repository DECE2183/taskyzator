package ui

import (
	"taskyzator/ui/model"
	mainpage "taskyzator/ui/model/mainPage"
)

func Run() {
	var err error

	err = mainpage.New().Run()
	if err != nil {
		model.PrettyExit(err, 6)
	}
}
