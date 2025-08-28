package reminder

import (
	"fmt"
	"time"

	"github.com/leksusdev/calendarOfEvents/datetime"
)

type Reminder struct {
	Message string       `json:"message"`
	At      time.Time    `json:"at"`
	Sent    bool         `json:"sent"`
	Timer   *time.Timer  `json:"-"`
	Notify  func(string) `json:"-"`
}

func NewReminder(message string, at time.Time, notify func(string)) (*Reminder, error) {
	msg, err := validateMessage(message)
	if err != nil {
		return nil, err
	}

	at = datetime.NormalizeUTCSeconds(at)

	if err := validateAt(at); err != nil {
		return nil, err
	}

	return &Reminder{
		Message: msg,
		At:      at,
		Sent:    false,
		Notify:  notify,
	}, nil
}

func (r *Reminder) Start() {
	if r.Timer != nil {
		r.Timer.Stop()
	}
	r.Timer = time.AfterFunc(time.Until(r.At), r.Send)
}

func (r *Reminder) Send() {
	if r.Sent {
		return
	}
	if r.Notify != nil {
		r.Notify(fmt.Sprintf("Напоминание: \"%s\" - \"%s\"", r.Message, datetime.FormatLocal(r.At)))
	}
	r.Sent = true
}

func (r *Reminder) Stop() {
	if r.Timer != nil {
		r.Timer.Stop()
		r.Timer = nil
	}
}
