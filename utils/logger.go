package utils

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

var (
	// Emoji constants
	SUCCESS = "✨"
	INFO    = "ℹ️"
	WARNING = "⚠️"
	ERROR   = "❌"
	DB      = "🗃️"
	AUTH    = "🔐"
	API     = "🌐"
	USER    = "👤"
	PROJECT = "📁"
	TASK    = "📝"
	COMMENT = "💬"

	// Color functions
	successColor = color.New(color.FgGreen, color.Bold).SprintFunc()
	infoColor    = color.New(color.FgCyan).SprintFunc()
	warningColor = color.New(color.FgYellow).SprintFunc()
	errorColor   = color.New(color.FgRed, color.Bold).SprintFunc()
	timeColor    = color.New(color.FgWhite).SprintFunc()
	methodColor  = color.New(color.FgMagenta, color.Bold).SprintFunc()
	pathColor    = color.New(color.FgBlue).SprintFunc()
	detailColor  = color.New(color.FgWhite).SprintFunc()
)

type Logger struct{}

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) LogRequest(method, path string, status int, latency time.Duration) {
	timestamp := time.Now().Format("15:04:05")
	
	var statusEmoji string
	var statusText string
	
	switch {
	case status >= 500:
		statusEmoji = ERROR
		statusText = errorColor(fmt.Sprintf("%d", status))
	case status >= 400:
		statusEmoji = WARNING
		statusText = warningColor(fmt.Sprintf("%d", status))
	case status >= 300:
		statusEmoji = INFO
		statusText = infoColor(fmt.Sprintf("%d", status))
	default:
		statusEmoji = SUCCESS
		statusText = successColor(fmt.Sprintf("%d", status))
	}

	// Get emoji based on path
	var contextEmoji string
	switch {
	case contains(path, "auth"):
		contextEmoji = AUTH
	case contains(path, "projects"):
		contextEmoji = PROJECT
	case contains(path, "tasks"):
		contextEmoji = TASK
	case contains(path, "comments"):
		contextEmoji = COMMENT
	default:
		contextEmoji = API
	}

	fmt.Printf("%s %s │ %s │ %s %s %s │ %s │ %s %s\n",
		contextEmoji,
		timeColor(timestamp),
		methodColor(fmt.Sprintf("%-6s", method)),
		statusEmoji,
		statusText,
		detailColor(fmt.Sprintf("(%dms)", latency.Milliseconds())),
		pathColor(path),
		USER,
		"::1",
	)
}

func (l *Logger) LogDB(operation, query string, rows int64, latency time.Duration) {
	timestamp := time.Now().Format("15:04:05")
	
	fmt.Printf("%s %s │ %s │ %s │ %s %s\n",
		DB,
		timeColor(timestamp),
		infoColor(fmt.Sprintf("%-8s", operation)),
		detailColor(fmt.Sprintf("Rows: %d", rows)),
		detailColor(fmt.Sprintf("(%dms)", latency.Milliseconds())),
		pathColor(query),
	)
}

func (l *Logger) LogInfo(context, message string) {
	timestamp := time.Now().Format("15:04:05")
	fmt.Printf("%s %s │ %s\n",
		INFO,
		timeColor(timestamp),
		infoColor(message),
	)
}

func (l *Logger) LogSuccess(context, message string) {
	timestamp := time.Now().Format("15:04:05")
	fmt.Printf("%s %s │ %s\n",
		SUCCESS,
		timeColor(timestamp),
		successColor(message),
	)
}

func (l *Logger) LogWarning(context, message string) {
	timestamp := time.Now().Format("15:04:05")
	fmt.Printf("%s %s │ %s\n",
		WARNING,
		timeColor(timestamp),
		warningColor(message),
	)
}

func (l *Logger) LogError(context, message string) {
	timestamp := time.Now().Format("15:04:05")
	fmt.Printf("%s %s │ %s\n",
		ERROR,
		timeColor(timestamp),
		errorColor(message),
	)
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[0:len(substr)] == substr
} 