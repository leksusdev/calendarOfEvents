package events

import (
	"strings"
	"testing"
)

func TestIsValidTitle(t *testing.T) {
	if isValidTitle("Ab") {
		t.Error("Ожидали false для слишком короткого заголовка, получили true")
	}

	long := strings.Repeat("A", 51)
	if isValidTitle(long) {
		t.Error("Ожидали false для слишком длинного заголовка, получили true")
	}

	if isValidTitle("Hi!;:*^$%") {
		t.Error("Ожидали false для заголовка с запрещёнными символами, получили true")
	}

	if !isValidTitle("Заголовок 42.,") {
		t.Error("Ожидали true для валидного заголовка, получили false")
	}
}
