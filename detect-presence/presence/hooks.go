package presence

import (
	"context"
	"log"
)

type OnLeaveHook interface {
	OnLeave(ctx context.Context, t *Tracker)
}

type OnReturnHook interface {
	OnReturn(ctx context.Context, t *Tracker)
}

func (t *Tracker) OnLeave(hook OnLeaveHook) {
	t.onLeaveHooks = append(t.onLeaveHooks, hook)
}

func (t *Tracker) OnReturn(hook OnReturnHook) {
	t.onReturnHooks = append(t.onReturnHooks, hook)
}

type OnLeaveFunc func(context.Context, *Tracker)
type OnReturnFunc func(context.Context, *Tracker)

func (fn OnLeaveFunc) OnLeave(ctx context.Context, t *Tracker) {
	fn(ctx, t)
}

func (fn OnReturnFunc) OnReturn(ctx context.Context, t *Tracker) {
	fn(ctx, t)
}

type loggingHook struct{}

func (loggingHook) OnLeave(context.Context, *Tracker) {
	log.Printf("Transitioned from home to away")
}

func (loggingHook) OnReturn(context.Context, *Tracker) {
	log.Printf("Transitioned from away to home")
}
