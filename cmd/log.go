package cmd

import (
	"strings"
	"sync"

	"github.com/leksusdev/calendarOfEvents/config"
	"github.com/leksusdev/calendarOfEvents/storage"
)

type LogHandler struct {
	mu      sync.Mutex
	lines   []string
	storage storage.Store
}

func NewLogHandler() *LogHandler {
	return &LogHandler{
		storage: storage.NewZipStorage(config.LogArchiveName),
	}
}

func (l *LogHandler) AddLine(s string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.lines = append(l.lines, s)
}

func (l *LogHandler) Output(s string) string {
	l.AddLine(s)
	return s
}

func (l *LogHandler) OutputLn(s string) string {
	return l.Output(s + "\n")
}

func (l *LogHandler) GetSnapshot() []string {
	l.mu.Lock()
	defer l.mu.Unlock()
	cp := make([]string, len(l.lines))
	copy(cp, l.lines)
	return cp
}

func (l *LogHandler) Save() error {
	snapshot := l.GetSnapshot()
	data := strings.Join(snapshot, "")
	return l.storage.Save([]byte(data))
}

func (l *LogHandler) Load() error {
	data, err := l.storage.Load()
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")

	l.mu.Lock()
	defer l.mu.Unlock()

	l.lines = make([]string, 0, len(lines))
	for _, line := range lines {
		if line == "" {
			continue
		}
		l.lines = append(l.lines, line+"\n")
	}
	return nil
}
