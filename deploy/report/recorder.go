package report

import (
	"fmt"
	"sync"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	deploypb "github.com/mjm/pi-tools/deploy/proto/deploy"
)

type Recorder struct {
	report deploypb.Report
	lock   sync.Mutex
}

func (r *Recorder) SetDeployID(deployID int64) {
	r.report.DeployId = deployID
}

func (r *Recorder) SetCommitInfo(commitSha string, commitMessage string) {
	r.report.CommitSha = commitSha
	r.report.CommitMessage = commitMessage
}

func (r *Recorder) Info(format string, v ...interface{}) Event {
	return r.addEvent(deploypb.ReportEvent_INFO, format, v...)
}

func (r *Recorder) Warning(format string, v ...interface{}) Event {
	return r.addEvent(deploypb.ReportEvent_WARNING, format, v...)
}

func (r *Recorder) Error(format string, v ...interface{}) Event {
	return r.addEvent(deploypb.ReportEvent_ERROR, format, v...)
}

func (r *Recorder) AppendEvent(evt *deploypb.ReportEvent) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.report.Events = append(r.report.Events, evt)
}

func (r *Recorder) addEvent(level deploypb.ReportEvent_Level, format string, v ...interface{}) Event {
	r.lock.Lock()
	defer r.lock.Unlock()

	evt := &deploypb.ReportEvent{
		Timestamp: timestamppb.Now(),
		Level:     level,
		Summary:   fmt.Sprintf(format, v...),
	}
	r.report.Events = append(r.report.Events, evt)
	return &event{evt: evt}
}

func (r *Recorder) Marshal() ([]byte, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	return proto.Marshal(&r.report)
}

type Event interface {
	WithDescription(string, ...interface{}) Event
	WithError(error) Event
}

type event struct {
	evt *deploypb.ReportEvent
}

func (e *event) WithDescription(format string, v ...interface{}) Event {
	e.evt.Description = fmt.Sprintf(format, v...)
	return e
}

func (e *event) WithError(err error) Event {
	return e.WithDescription("Error: %v", err)
}
