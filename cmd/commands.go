package cmd

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/leksusdev/calendarOfEvents/config"
	"github.com/leksusdev/calendarOfEvents/datetime"
	"github.com/leksusdev/calendarOfEvents/events"
	"github.com/leksusdev/calendarOfEvents/logger"
)

const (
	addFormat          = "add <\"название события\"> <\"дата и время\"> <приоритет>"
	removeFormat       = "remove <ID>"
	updateFormat       = "update <ID> <\"название события\"> <\"дата и время\"> <приоритет>"
	remindFormat       = "remind <ID> <\"сообщение\"> <\"дата и время\"|duration>"
	cancelRemindFormat = "remind-cancel <ID>"
)

func (c *Cmd) handleAdd(parts []string) {
	logger.Info("Обработка команды add")
	if len(parts) < 4 {
		c.outputLn("Формат: " + addFormat)
		logger.Error("Неверный формат команды add")
		return
	}

	title := parts[1]
	date := parts[2]
	priority := events.Priority(parts[3])

	e, err := c.calendar.AddEvent(title, date, priority)
	if err != nil {
		c.outputLn("Ошибка: " + err.Error())
		logger.Error("Ошибка добавления события: " + err.Error())
		return
	}
	c.outputLn("Событие: \"" + e.Title + "\" добавлено")
	logger.Info(fmt.Sprintf("Событие добавлено: ID=%s, Title=%s", e.ID, e.Title))
}

func (c *Cmd) handleList() {
	logger.Info("Обработка команды list")
	eventsList := c.calendar.GetEvents()
	if len(eventsList) == 0 {
		c.outputLn("Календарь пуст")
		logger.Info("Календарь пуст")
		return
	}

	c.outputLn(fmt.Sprintf("|%-*s|%-*s|%-*s|%-*s",
		config.ListColWidthID, "ID:",
		config.ListColWidthTitle, "Событие:",
		config.ListColWidthDate, "Дата-время:",
		config.ListColWidthStatus, "Статус:"))
	for _, e := range eventsList {
		c.outputLn(fmt.Sprintf("|%-*s|%-*s|%-*s|%-*s",
			config.ListColWidthID, e.ID,
			config.ListColWidthTitle, e.Title+strings.Repeat(".", config.ListTitlePad-utf8.RuneCountInString(e.Title)),
			config.ListColWidthDate, datetime.FormatLocal(e.StartAt),
			config.ListColWidthStatus, e.Priority,
		))
	}
	logger.Info(fmt.Sprintf("Выведено %d событий", len(eventsList)))
}

func (c *Cmd) handleRemove(parts []string) {
	logger.Info("Обработка команды remove")
	if len(parts) < 2 {
		c.outputLn("Формат: " + removeFormat)
		logger.Error("Неверный формат команды remove")
		return
	}

	ID := parts[1]
	deletedEvent, err := c.calendar.DeleteEvent(ID)
	if err != nil {
		c.outputLn("Ошибка: " + err.Error())
		logger.Error("Ошибка удаления события: " + err.Error())
		return
	}
	c.outputLn("Событие удалено: \"" + deletedEvent.Title + "\"")
	logger.Info(fmt.Sprintf("Событие удалено: ID=%s, Title=%s", deletedEvent.ID, deletedEvent.Title))
}

func (c *Cmd) handleUpdate(parts []string) {
	logger.Info("Обработка команды update")
	if len(parts) < 5 {
		c.outputLn("Формат: " + updateFormat)
		logger.Error("Неверный формат команды update")
		return
	}

	ID := parts[1]
	newTitle := parts[2]
	newDate := parts[3]
	newPriority := events.Priority(parts[4])

	oldTitle, newTitle, err := c.calendar.EditEvent(ID, newTitle, newDate, newPriority)
	if err != nil {
		c.outputLn("Ошибка: " + err.Error())
		logger.Error("Ошибка обновления события: " + err.Error())
		return
	}
	c.outputLn(fmt.Sprintf("Событие обновлено: \"%s\" на \"%s\"", oldTitle, newTitle))
	logger.Info(fmt.Sprintf("Событие обновлено: ID=%s, OldTitle=%s, NewTitle=%s", ID, oldTitle, newTitle))
}

func (c *Cmd) handleRemind(parts []string) {
	logger.Info("Обработка команды remind")
	if len(parts) < 4 {
		c.outputLn("Формат: " + remindFormat)
		logger.Error("Неверный формат команды remind")
		return
	}

	id := parts[1]
	message := strings.TrimSpace(parts[2])
	at := parts[3]

	if err := c.calendar.SetEventReminder(id, message, at); err != nil {
		c.outputLn("Ошибка: " + err.Error())
		logger.Error("Ошибка добавления напоминания: " + err.Error())
		return
	}
	c.outputLn("Добавлено напоминание: \"" + message + "\"")
	logger.Info(fmt.Sprintf("Добавлено напоминание: ID=%s, Message=%s", id, message))
}

func (c *Cmd) handleRemindCancel(parts []string) {
	logger.Info("Обработка команды remind-cancel")
	if len(parts) < 2 {
		c.outputLn("Формат: " + cancelRemindFormat)
		logger.Error("Неверный формат команды remind-cancel")
		return
	}

	id := parts[1]
	if err := c.calendar.CancelEventReminder(id); err != nil {
		c.outputLn("Ошибка: " + err.Error())
		logger.Error("Ошибка отмены напоминания: " + err.Error())
		return
	}
	c.outputLn("Напоминание отменено")
	logger.Info(fmt.Sprintf("Напоминание отменено: ID=%s", id))
}

func (c *Cmd) handleHelp() {
	logger.Info("Обработка команды help")
	c.outputLn(".......................................................................................")
	c.outputLn(fmt.Sprintf(":             Добавить: %s         :", addFormat))
	c.outputLn(fmt.Sprintf(":              Удалить: %s                                                   :", removeFormat))
	c.outputLn(fmt.Sprintf(":             Обновить: %s :", updateFormat))
	c.outputLn(fmt.Sprintf(":          Напоминание: %s           :", remindFormat))
	c.outputLn(fmt.Sprintf(": Отменить напоминание: %s                                            :", cancelRemindFormat))
	c.outputLn(":               Список: list                                                          :")
	c.outputLn(":                  Лог: log                                                           :")
	c.outputLn(":        Сохранить лог: log-save                                                      :")
	c.outputLn(":        Загрузить лог: log-load                                                      :")
	c.outputLn(":                Выход: exit                                                          :")
	c.outputLn(":.....................................................................................:")
	c.outputLn(fmt.Sprintf(": Пример шаблона даты и времени: %s                                     :", datetime.LayoutFormat))
	c.outputLn(": Пример шаблона duration для напоминания: 1h50m30s                                   :")
	c.outputLn(fmt.Sprintf(": Допустимые приоритеты: %s, %s, %s                                            :", events.PriorityLow, events.PriorityMedium, events.PriorityHigh))
	c.outputLn(fmt.Sprintf(": Данные сохраняются в файл %s при выходе из программы                :", config.DataFileName))
	c.outputLn(fmt.Sprintf(": Логи команд сохраняются в файл %s и архивируются в %s    :", config.ZipLogEntryName, config.LogArchiveName))
	c.outputLn(fmt.Sprintf(": Логи приложения хранятся в файле %s                                       :", config.LogFileName))
	c.outputLn(": При обновлении события некоторые поля можно пропустить вводом символа <_>           :")
	c.outputLn(":.....................................................................................:")
}

func (c *Cmd) handleLog() {
	logger.Info("Обработка команды log")
	lines := c.logHandler.GetSnapshot()
	if len(lines) == 0 {
		c.outputLn("Лог пуст")
		logger.Info("Лог пуст")
		return
	}
	for _, line := range lines {
		c.output(line)
	}
	logger.Info(fmt.Sprintf("Выведено %d строк лога", len(lines)))
}

func (c *Cmd) handleLogSave() {
	logger.Info("Обработка команды log-save")
	if err := c.logHandler.Save(); err != nil {
		c.outputLn("Ошибка сохранения лога: " + err.Error())
		logger.Error("Ошибка сохранения лога: " + err.Error())
		return
	}
	c.outputLn("Лог сохранён")
	logger.Info("Лог сохранён")
}

func (c *Cmd) handleLogLoad() {
	logger.Info("Обработка команды log-load")
	if err := c.logHandler.Load(); err != nil {
		c.outputLn("Ошибка загрузки лога: " + err.Error())
		logger.Error("Ошибка загрузки лога: " + err.Error())
		return
	}
	c.outputLn("Лог загружен")
	logger.Info("Лог загружен")
}

func (c *Cmd) handleExit() {
	logger.Info("Обработка команды exit")
	for _, e := range c.calendar.GetEvents() {
		if e.Reminder != nil {
			e.Reminder.Stop()
		}
	}

	err := c.calendar.Save()
	if err != nil {
		c.outputLn("Ошибка сохранения данных: " + err.Error())
		logger.Error("Ошибка сохранения данных: " + err.Error())
		return
	}
	c.calendar.Close()
	c.wg.Wait()
	logger.Info("Приложение завершило работу")
	os.Exit(0)
}
