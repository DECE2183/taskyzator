package ui

import (
	"github.com/dece2183/taskyzator/ui/model"
	mainpage "github.com/dece2183/taskyzator/ui/model/mainPage"
)

func Run() {
	var err error

	err = mainpage.New().Run()
	if err != nil {
		model.PrettyExit(err, 6)
	}
}
