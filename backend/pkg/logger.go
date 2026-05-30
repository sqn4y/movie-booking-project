package pkg

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"sync"
)

const (
	colorReset     = "\033[0m"
	colorRed       = "\033[31m"
	colorGreen     = "\033[32m"
	colorYellow    = "\033[33m"
	colorBlue      = "\033[34m"
	colorMagenta   = "\033[35m"
	colorCyan      = "\033[36m"
	colorWhite     = "\033[37m"
	colorLightCyan = "\033[96m" // светлоголубой для имени логгера
)

type springBootHandler struct {
	out     io.Writer
	level   slog.Leveler
	appName string
	pid     int
	attrs   []slog.Attr
	mu      *sync.Mutex
}

func NewLogger(out io.Writer, level slog.Leveler, appName string) *slog.Logger {
	return slog.New(&springBootHandler{
		out:     out,
		level:   level,
		appName: appName,
		pid:     os.Getpid(),
		mu:      &sync.Mutex{},
	})
}

func (h *springBootHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level.Level()
}

func (h *springBootHandler) Handle(_ context.Context, record slog.Record) error {
	var builder strings.Builder

	// 1. Время с миллисекундами и таймзоной (ISO 8601)
	builder.WriteString(record.Time.Format("2006-01-02T15:04:05.000Z07:00"))
	builder.WriteByte(' ')

	// 2. Уровень с цветом и фиксированной шириной (5 символов)
	levelName := h.levelName(record.Level)
	builder.WriteString(h.levelColor(record.Level))
	builder.WriteString(fmt.Sprintf("%-5s", levelName))
	builder.WriteString(colorReset)
	builder.WriteByte(' ')

	// 3. PID (фиксированная ширина 5 символов)
	builder.WriteString(fmt.Sprintf("%5d", h.pid))
	builder.WriteByte(' ')

	// 4. Имя приложения
	builder.WriteString("--- [")
	builder.WriteString(h.appName)
	builder.WriteString("]")

	// 5. Имя потока/горутины (фиксированная ширина 10 символов)
	builder.WriteString(" [")
	builder.WriteString(h.threadName())
	builder.WriteString("] ")

	// 6. Имя логгера (пакет/функция) - светлоголубым
	loggerName := h.getLoggerName(record)
	// Фиксированная ширина 50 символов для выравнивания сообщения
	paddedLoggerName := fmt.Sprintf("%-50s", loggerName)
	builder.WriteString(colorLightCyan)
	builder.WriteString(paddedLoggerName)
	builder.WriteString(colorReset)
	builder.WriteString(" : ")

	// 7. Сообщение с подсветкой HTTP-элементов
	message := record.Message
	message = h.highlightHTTP(message)
	builder.WriteString(message)

	// 8. Дополнительные атрибуты
	record.Attrs(func(attr slog.Attr) bool {
		builder.WriteString(", ")
		builder.WriteString(attr.Key)
		builder.WriteString("=")
		// Подсвечиваем значения атрибутов
		value := h.highlightValue(attr.Key, attr.Value.Any())
		builder.WriteString(value)
		return true
	})

	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := fmt.Fprintln(h.out, builder.String())
	return err
}

func (h *springBootHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	next := h.clone()
	next.attrs = append(next.attrs, attrs...)
	return next
}

func (h *springBootHandler) WithGroup(name string) slog.Handler {
	return h
}

func (h *springBootHandler) clone() *springBootHandler {
	next := *h
	next.attrs = append([]slog.Attr(nil), h.attrs...)
	return &next
}

// highlightHTTP подсвечивает HTTP-методы, пути и статусы
func (h *springBootHandler) highlightHTTP(msg string) string {
	// HTTP методы
	methods := map[string]string{
		"GET":     colorCyan,
		"POST":    colorGreen,
		"PUT":     colorBlue,
		"DELETE":  colorRed,
		"PATCH":   colorYellow,
		"OPTIONS": colorMagenta,
		"HEAD":    colorMagenta,
	}

	for method, color := range methods {
		if strings.Contains(msg, method) {
			msg = strings.ReplaceAll(msg, method, color+method+colorReset)
		}
	}

	// HTTP статусы (2xx, 3xx, 4xx, 5xx)
	// Ищем паттерны типа "200", "404", "500" и т.д.
	words := strings.Fields(msg)
	for i, word := range words {
		// Проверяем, похоже ли слово на статус (3 цифры)
		if len(word) == 3 {
			// Проверяем, что все символы - цифры
			isDigit := true
			for _, c := range word {
				if c < '0' || c > '9' {
					isDigit = false
					break
				}
			}
			if isDigit {
				statusCode := word
				statusColor := h.getStatusColor(statusCode)
				words[i] = statusColor + statusCode + colorReset
			}
		}
	}

	return strings.Join(words, " ")
}

// highlightValue подсвечивает значения атрибутов
func (h *springBootHandler) highlightValue(key string, value interface{}) string {
	strValue := fmt.Sprintf("%v", value)

	// Подсветка HTTP-путей
	if key == "path" || key == "uri" || key == "url" {
		return colorCyan + strValue + colorReset
	}

	// Подсветка статусов
	if key == "status" || key == "status_code" {
		return h.getStatusColor(strValue) + strValue + colorReset
	}

	// Подсветка HTTP-методов
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"}
	for _, method := range methods {
		if strValue == method {
			switch method {
			case "GET":
				return colorCyan + strValue + colorReset
			case "POST":
				return colorGreen + strValue + colorReset
			case "PUT":
				return colorBlue + strValue + colorReset
			case "DELETE":
				return colorRed + strValue + colorReset
			default:
				return colorYellow + strValue + colorReset
			}
		}
	}

	// Подсветка чисел (ID)
	if key == "id" || key == "movie_id" || key == "user_id" {
		return colorYellow + strValue + colorReset
	}

	return strValue
}

func (h *springBootHandler) getStatusColor(status string) string {
	if len(status) != 3 {
		return colorWhite
	}

	switch status[0] {
	case '1':
		return colorCyan // 1xx - informational
	case '2':
		return colorGreen // 2xx - success
	case '3':
		return colorBlue // 3xx - redirect
	case '4':
		return colorYellow // 4xx - client error
	case '5':
		return colorRed // 5xx - server error
	default:
		return colorWhite
	}
}

func (h *springBootHandler) levelName(level slog.Level) string {
	switch {
	case level >= slog.LevelError:
		return "ERROR"
	case level >= slog.LevelWarn:
		return "WARN"
	case level >= slog.LevelInfo:
		return "INFO"
	default:
		return "DEBUG"
	}
}

func (h *springBootHandler) levelColor(level slog.Level) string {
	switch {
	case level >= slog.LevelError:
		return colorRed
	case level >= slog.LevelWarn:
		return colorYellow
	case level >= slog.LevelInfo:
		return colorGreen
	default:
		return colorCyan
	}
}

func (h *springBootHandler) threadName() string {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	// Получаем ID горутины из стека
	id := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	// Делаем ID фиксированной ширины (10 символов)
	return fmt.Sprintf("%-10s", id)
}

func (h *springBootHandler) getLoggerName(record slog.Record) string {
	if record.PC == 0 {
		return "unknown"
	}

	frames := runtime.CallersFrames([]uintptr{record.PC})
	frame, _ := frames.Next()

	if frame.Function == "" {
		return "unknown"
	}

	// Оставляем только имя функции без полного пути
	parts := strings.Split(frame.Function, "/")
	if len(parts) > 0 {
		lastPart := parts[len(parts)-1]
		return lastPart
	}
	return frame.Function
}

// WithThread - добавляет имя потока в контекст (опционально)
func WithThread(ctx context.Context, threadName string) context.Context {
	return context.WithValue(ctx, "thread", threadName)
}
