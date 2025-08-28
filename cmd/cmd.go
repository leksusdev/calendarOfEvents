package cmd

import (
	"strings"
	"sync"

	"github.com/c-bata/go-prompt"
	"github.com/google/shlex"
	"github.com/leksusdev/calendarOfEvents/calendar"
	"github.com/leksusdev/calendarOfEvents/config"
)

type Cmd struct {
	calendar   *calendar.Calendar
	wg         sync.WaitGroup
	logHandler *LogHandler
}

func NewCmd(c *calendar.Calendar) *Cmd {
	return &Cmd{
		calendar:   c,
		logHandler: NewLogHandler(),
	}
}

func (c *Cmd) output(s string) {
	print(s)
	c.logHandler.AddLine(s)
}

func (c *Cmd) outputLn(s string) {
	c.output(s + "\n")
}

func (c *Cmd) executor(input string) {
	input = strings.TrimSpace(input)
	if input == "" {
		return
	}
	c.logHandler.AddLine(config.PromptPrefix + input + "\n")

	parts, err := shlex.Split(input)
	if err != nil {
		c.outputLn(err.Error())
		return
	}
	if len(parts) == 0 {
		return
	}

	cmd := strings.ToLower(parts[0])

	switch cmd {
	case "add":
		c.handleAdd(parts)
	case "list":
		c.handleList()
	case "remove":
		c.handleRemove(parts)
	case "update":
		c.handleUpdate(parts)
	case "remind":
		c.handleRemind(parts)
	case "remind-cancel":
		c.handleRemindCancel(parts)
	case "help":
		c.handleHelp()
	case "log":
		c.handleLog()
	case "log-save":
		c.handleLogSave()
	case "log-load":
		c.handleLogLoad()
	case "exit":
		c.handleExit()
	default:
		c.outputLn("Неизвестная команда")
		c.outputLn("Введите <help> для информации")
	}
}

func (c *Cmd) completer(d prompt.Document) []prompt.Suggest {
	if strings.Contains(d.TextBeforeCursor(), " ") {
		return []prompt.Suggest{}
	}
	suggestions := []prompt.Suggest{
		{Text: "add", Description: "Добавить событие"},
		{Text: "list", Description: "Показать все события"},
		{Text: "update", Description: "Обновить событие"},
		{Text: "remove", Description: "Удалить событие"},
		{Text: "remind", Description: "Добавить напоминание к событию"},
		{Text: "remind-cancel", Description: "Отменить напоминание к событию"},
		{Text: "help", Description: "Описание команд"},
		{Text: "log", Description: "Показать лог сессии"},
		{Text: "log-save", Description: "Сохранить лог в файл"},
		{Text: "log-load", Description: "Загрузить лог из файла"},
		{Text: "exit", Description: "Выйти из программы"},
	}

	return prompt.FilterHasPrefix(suggestions, d.GetWordBeforeCursor(), true)
}

func (c *Cmd) Run() {
	p := prompt.New(
		c.executor,
		c.completer,
		prompt.OptionPrefix(config.PromptPrefix),
		prompt.OptionMaxSuggestion(config.PromptMaxSuggestions),
	)
	c.wg.Go(func() {
		for msg := range c.calendar.Notification {
			c.outputLn(msg)
		}
	})
	p.Run()
}
