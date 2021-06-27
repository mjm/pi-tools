package nomadic

import (
	"context"
	"fmt"
	"log"

	deploypb "github.com/mjm/pi-tools/deploy/proto/deploy"
	"github.com/mjm/pi-tools/deploy/report"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	Info(format string, v ...interface{})
	Warning(format string, v ...interface{})
	Error(format string, v ...interface{})
}

type EventOption func(e report.Event)

func WithError(err error) EventOption {
	return func(e report.Event) {
		e.WithError(err)
	}
}

func withDescription(format string, v ...interface{}) EventOption {
	return func(e report.Event) {
		e.WithDescription(format, v...)
	}
}

func EventArgsAndOptions(v ...interface{}) (args []interface{}, opts []EventOption) {
	for _, x := range v {
		if x, ok := x.(EventOption); ok {
			opts = append(opts, x)
		} else {
			args = append(args, x)
		}
	}
	return
}

type logEvents struct{}

func NewLoggingEventReporter() EventReporter {
	return logEvents{}
}

func (logEvents) Info(format string, v ...interface{}) {
	args, opts := EventArgsAndOptions(v)
	log.Printf(format, args...)
	e := logEventImpl{}
	for _, opt := range opts {
		opt(e)
	}
}

func (logEvents) Warning(format string, v ...interface{}) {
	args, opts := EventArgsAndOptions(v)
	log.Printf("warning: "+format, args...)
	e := logEventImpl{}
	for _, opt := range opts {
		opt(e)
	}
}

func (logEvents) Error(format string, v ...interface{}) {
	args, opts := EventArgsAndOptions(v)
	log.Printf("error: "+format, args...)
	e := logEventImpl{}
	for _, opt := range opts {
		opt(e)
	}
}

type logEventImpl struct{}

func (l logEventImpl) WithDescription(format string, v ...interface{}) report.Event {
	log.Printf("  "+format, v...)
	return l
}

func (l logEventImpl) WithError(err error) report.Event {
	return l.WithDescription("%s", err.Error())
}

type channelEvents struct {
	ch chan<- *deploypb.ReportEvent
}

func NewChannelEventReporter(ch chan<- *deploypb.ReportEvent) EventReporter {
	return &channelEvents{ch: ch}
}

func (c *channelEvents) Info(format string, v ...interface{}) {
	c.emitEvent(deploypb.ReportEvent_INFO, format, v...)
}

func (c *channelEvents) Warning(format string, v ...interface{}) {
	c.emitEvent(deploypb.ReportEvent_WARNING, format, v...)
}

func (c *channelEvents) Error(format string, v ...interface{}) {
	c.emitEvent(deploypb.ReportEvent_ERROR, format, v...)
}

func (c *channelEvents) emitEvent(level deploypb.ReportEvent_Level, format string, v ...interface{}) {
	args, opts := EventArgsAndOptions(v)
	summary := fmt.Sprintf(format, args...)
	evt := wrappedEvent{
		Timestamp: timestamppb.Now(),
		Level:     level,
		Summary:   summary,
	}
	for _, opt := range opts {
		opt(&evt)
	}
	e := (*deploypb.ReportEvent)(&evt)
	c.ch <- e
}

type wrappedEvent deploypb.ReportEvent

func (e *wrappedEvent) WithDescription(format string, v ...interface{}) report.Event {
	e.Description = fmt.Sprintf(format, v...)
	return e
}

func (e *wrappedEvent) WithError(err error) report.Event {
	return e.WithDescription("Error: %v", err)
}
