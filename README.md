# Calendar of Events

Консольное приложение календарь событий с возможностью установки напоминаний.

## Сборка

```bash
# macOS (arm)
GOOS=darwin GOARCH=arm64 go build -o calendar-darwin-arm64

# Windows (x64)
GOOS=windows GOARCH=amd64 go build -o calendar-windows-amd64.exe
```