package events

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/leksusdev/calendarOfEvents/datetime"
	"github.com/leksusdev/calendarOfEvents/logger"
	"github.com/leksusdev/calendarOfEvents/reminder"
)

type Event struct {
	ID       string             `json:"id"`
	Title    string             `json:"title"`
	StartAt  time.Time          `json:"start_at"`
	Priority Priority           `json:"priority"`
	Reminder *reminder.Reminder `json:"reminder"`
}

var (
	ErrInvalidTitle      = errors.New("неверный формат заголовка")
	ErrInvalidDate       = errors.New("неверный формат даты")
	ErrEmptyReminderTime = errors.New("время напоминания не может быть пустым")
	ErrZeroDuration      = errors.New("время должно быть больше нуля")
)

func makeEvent(id string, title string, dateStr string, p Priority, reminder *reminder.Reminder) (Event, error) {
	title = strings.TrimSpace(title)
	dateStr = strings.TrimSpace(dateStr)

	if !isValidTitle(title) {
		return Event{}, fmt.Errorf("ошибка проверки заголовка: %w", ErrInvalidTitle)
	}

	t, err := datetime.ParseLocal(dateStr)
	if err != nil {
		return Event{}, fmt.Errorf("ошибка проверки даты/времени: %w", ErrInvalidDate)
	}

	t = datetime.NormalizeUTCSeconds(t)

	if err := p.Validate(); err != nil {
		return Event{}, err
	}
	return Event{
		ID:       id,
		Title:    title,
		StartAt:  t,
		Priority: p,
		Reminder: reminder,
	}, nil
}

func NewEvent(title string, dateStr string, priority Priority) (*Event, error) {
	event, err := makeEvent(generateUUID(), title, dateStr, priority, nil)
	if err != nil {
		logger.Error(fmt.Sprintf("Ошибка создания события: %v", err))
		return nil, err
	}
	logger.Info(fmt.Sprintf("Создано событие: ID=%s, Title=%s", event.ID, event.Title))
	return &event, nil
}

func (e *Event) Update(title string, dateStr string, priority Priority) error {
	updatedEvent, err := makeEvent(e.ID, title, dateStr, priority, e.Reminder)
	if err != nil {
		logger.Error(fmt.Sprintf("Ошибка обновления события ID=%s: %v", e.ID, err))
		return err
	}
	e.Title = updatedEvent.Title
	e.StartAt = updatedEvent.StartAt
	e.Priority = updatedEvent.Priority
	logger.Info(fmt.Sprintf("Обновлено событие: ID=%s, NewTitle=%s", e.ID, e.Title))
	return nil
}

func (e *Event) AddReminder(message string, at string, notify func(string)) error {
	at = strings.TrimSpace(at)
	if at == "" {
		err := fmt.Errorf("ошибка проверки даты/времени: %w", ErrEmptyReminderTime)
		logger.Error(fmt.Sprintf("Ошибка добавления напоминания для события ID=%s: %v", e.ID, err))
		return err
	}

	var t time.Time

	if d, err := time.ParseDuration(at); err == nil {
		if d <= 0 {
			err := fmt.Errorf("ошибка проверки даты/времени: %w", ErrZeroDuration)
			logger.Error(fmt.Sprintf("Ошибка добавления напоминания для события ID=%s: %v", e.ID, err))
			return err
		}
		t = time.Now().Add(d)
	} else {
		tt, err2 := datetime.ParseLocal(at)
		if err2 != nil {
			err := fmt.Errorf("ошибка проверки даты/времени: %w", ErrInvalidDate)
			logger.Error(fmt.Sprintf("Ошибка добавления напоминания для события ID=%s: %v", e.ID, err))
			return err
		}
		t = tt
	}

	t = datetime.NormalizeUTCSeconds(t)

	r, err := reminder.NewReminder(message, t, notify)
	if err != nil {
		logger.Error(fmt.Sprintf("Ошибка создания напоминания для события ID=%s: %v", e.ID, err))
		return err
	}
	e.Reminder = r
	e.Reminder.Start()
	logger.Info(fmt.Sprintf("Добавлено напоминание для события ID=%s: Message=%s", e.ID, message))
	return nil
}

func (e *Event) RemoveReminder() {
	if e.Reminder != nil {
		e.Reminder.Stop()
		logger.Info(fmt.Sprintf("Напоминание остановлено для события ID=%s", e.ID))
		e.Reminder = nil
	}
}
