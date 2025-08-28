package reminder

import (
	"errors"
	"strings"
	"testing"
)

func TestValidateMessage(t *testing.T) {
	_, err := validateMessage(" ")
	if !errors.Is(err, ErrEmptyMessage) {
		t.Errorf("Ожидали ErrEmptyMessage, получили: %v", err)
	}

	long := strings.Repeat("A", 101)
	_, err = validateMessage(long)
	if !errors.Is(err, ErrMessageTooLong) {
		t.Errorf("Ожидали ErrMessageTooLong, получили: %v", err)
	}

	_, err = validateMessage("Сообщение!")
	if err != nil {
		t.Errorf("Не ожидали ошибку для корректного сообщения, получили: %v", err)
	}
}
