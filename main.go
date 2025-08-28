package main

import (
	"fmt"

	"github.com/leksusdev/calendarOfEvents/calendar"
	"github.com/leksusdev/calendarOfEvents/cmd"
	"github.com/leksusdev/calendarOfEvents/config"
	"github.com/leksusdev/calendarOfEvents/logger"
	"github.com/leksusdev/calendarOfEvents/storage"
)

func main() {
	if err := logger.Init(config.LogFileName); err != nil {
		fmt.Printf("Ошибка инициализации логгера: %v\n", err)
		return
	}
	defer logger.Close()

	logger.Info("Приложение запущено")

	s := storage.NewJsonStorage(config.DataFileName)
	c := calendar.NewCalendar(s)
	if err := c.Load(); err != nil {
		logger.Error(fmt.Sprintf("Ошибка загрузки данных: %s", err))
		return
	}

	cli := cmd.NewCmd(c)
	logger.Info("Запуск командной оболочки")
	cli.Run()
}
