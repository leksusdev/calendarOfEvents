package reminder

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const maxMessageLen = 100

var (
	ErrEmptyMessage   = errors.New("сообщение не может быть пустым")
	ErrMessageTooLong = errors.New("сообщение слишком длинное")
	ErrZeroTime       = errors.New("время не может быть нулевым")
	ErrPastTime       = errors.New("время не может быть в прошлом")
)

func validateMessage(msg string) (string, error) {
	msg = strings.TrimSpace(msg)
	if msg == "" {
		return "", fmt.Errorf("ошибка проверки сообщения: %w", ErrEmptyMessage)
	}
	if len([]rune(msg)) > maxMessageLen {
		return "", fmt.Errorf("ошибка проверки сообщения: %w", ErrMessageTooLong)
	}
	return msg, nil
}

func validateAt(at time.Time) error {
	if at.IsZero() {
		return fmt.Errorf("ошибка проверки даты/времени: %w", ErrZeroTime)
	}
	if at.Before(time.Now()) {
		return fmt.Errorf("ошибка проверки даты/времени: %w", ErrPastTime)
	}
	return nil
}
