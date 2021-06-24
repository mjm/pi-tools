package nomadic

import (
	"context"
	"log"

	"github.com/mjm/pi-tools/deploy/report"
)

type contextKey string

var contextKeyEvents = contextKey("events")

func WithEvents(ctx context.Context, events EventReporter) context.Context {
	return context.WithValue(ctx, contextKeyEvents, events)
}

func Events(ctx context.Context) EventReporter {
	events, _ := ctx.Value(contextKeyEvents).(EventReporter)
	return events
}

type EventReporter interface {
	Info(format string, v ...interface{}) report.Event
	Warning(format string, v ...interface{}) report.Event
	Error(format string, v ...interface{}) report.Event
}

type logEvents struct{}

func NewLoggingEventReporter() EventReporter {
	return logEvents{}
}

func (logEvents) Info(format string, v ...interface{}) report.Event {
	log.Printf(format, v...)
	return logEventImpl{}
}

func (logEvents) Warning(format string, v ...interface{}) report.Event {
	log.Printf("warning: " + format, v...)
	return logEventImpl{}
}

func (logEvents) Error(format string, v ...interface{}) report.Event {
	log.Printf("error: " + format, v...)
	return logEventImpl{}
}

type logEventImpl struct{}

func (l logEventImpl) WithDescription(format string, v ...interface{}) report.Event {
	log.Printf("  " + format, v...)
	return l
}

func (l logEventImpl) WithError(err error) report.Event {
	return l.WithDescription("%s", err.Error())
}
