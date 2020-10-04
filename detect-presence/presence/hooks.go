package presence

import (
	"log"
)

type OnLeaveHook interface {
	OnLeave(t *Tracker)
}

type OnReturnHook interface {
	OnReturn(t *Tracker)
}

func (t *Tracker) OnLeave(hook OnLeaveHook) {
	t.onLeaveHooks = append(t.onLeaveHooks, hook)
}

func (t *Tracker) OnReturn(hook OnReturnHook) {
	t.onReturnHooks = append(t.onReturnHooks, hook)
}

type OnLeaveFunc func(*Tracker)
type OnReturnFunc func(*Tracker)

func (fn OnLeaveFunc) OnLeave(t *Tracker) {
	fn(t)
}

func (fn OnReturnFunc) OnReturn(t *Tracker) {
	fn(t)
}

type loggingHook struct{}

func (loggingHook) OnLeave(_ *Tracker) {
	log.Printf("Transitioned from home to away")
}

func (loggingHook) OnReturn(_ *Tracker) {
	log.Printf("Transitioned from away to home")
}
