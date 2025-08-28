package calendar

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/leksusdev/calendarOfEvents/config"
	"github.com/leksusdev/calendarOfEvents/datetime"
	"github.com/leksusdev/calendarOfEvents/events"
	"github.com/leksusdev/calendarOfEvents/storage"
)

type Calendar struct {
	calendarEvents map[string]*events.Event
	storage        storage.Store
	Notification   chan string
}

var (
	ErrEventNotFound    = errors.New("событие не найдено")
	ErrReminderNotFound = errors.New("у события нет напоминания")
)

func NewCalendar(s storage.Store) *Calendar {
	return &Calendar{
		calendarEvents: make(map[string]*events.Event),
		storage:        s,
		Notification:   make(chan string),
	}
}

func (c *Calendar) Save() error {
	var data []byte
	var err error
	if config.PrettyJSON {
		data, err = json.MarshalIndent(c.calendarEvents, "", "  ")
	} else {
		data, err = json.Marshal(c.calendarEvents)
	}
	if err != nil {
		return fmt.Errorf("ошибка сериализации JSON: %w", err)
	}
	return c.storage.Save(data)
}

func (c *Calendar) Load() error {
	data, err := c.storage.Load()
	if err != nil {
		return fmt.Errorf("ошибка загрузки из стораджа: %w", err)
	}

	if len(data) == 0 {
		return nil
	}

	if err := json.Unmarshal(data, &c.calendarEvents); err != nil {
		return fmt.Errorf("ошибка парсинга JSON: %w", err)
	}
	return nil
}

func (c *Calendar) AddEvent(title string, dateStr string, priority events.Priority) (*events.Event, error) {
	e, err := events.NewEvent(title, dateStr, priority)
	if err != nil {
		return nil, err
	}

	c.calendarEvents[e.ID] = e
	return e, nil
}

func (c *Calendar) GetEvents() []*events.Event {
	eventsList := make([]*events.Event, 0, len(c.calendarEvents))
	for _, e := range c.calendarEvents {
		eventsList = append(eventsList, e)
	}
	return eventsList
}

func (c *Calendar) DeleteEvent(id string) (*events.Event, error) {
	e, exists := c.calendarEvents[id]
	if !exists {
		return nil, fmt.Errorf("id=%q: %w", id, ErrEventNotFound)
	}
	if e.Reminder != nil {
		e.RemoveReminder()
	}
	delete(c.calendarEvents, id)
	return e, nil
}

func (c *Calendar) EditEvent(id string, title string, dateStr string, priority events.Priority) (string, string, error) {
	e, exists := c.calendarEvents[id]
	if !exists {
		return "", "", fmt.Errorf("id=%q: %w", id, ErrEventNotFound)
	}

	oldTitle := e.Title

	if title == "_" {
		title = e.Title
	}
	if dateStr == "_" {
		dateStr = datetime.FormatLocal(e.StartAt)
	}
	if priority == "_" {
		priority = e.Priority
	}

	err := e.Update(title, dateStr, priority)
	if err != nil {
		return "", "", err
	}
	return oldTitle, e.Title, nil
}

func (c *Calendar) Notify(msg string) {
	c.Notification <- msg
}

func (c *Calendar) Close() {
	close(c.Notification)
}

func (c *Calendar) SetEventReminder(id string, message string, at string) error {
	e, exists := c.calendarEvents[id]
	if !exists {
		return fmt.Errorf("id=%q: %w", id, ErrEventNotFound)
	}
	return e.AddReminder(message, at, c.Notify)
}

func (c *Calendar) CancelEventReminder(id string) error {
	e, exists := c.calendarEvents[id]
	if !exists {
		return fmt.Errorf("id=%q: %w", id, ErrEventNotFound)
	}

	if e.Reminder == nil {
		return ErrReminderNotFound
	}

	e.RemoveReminder()
	return nil
}
