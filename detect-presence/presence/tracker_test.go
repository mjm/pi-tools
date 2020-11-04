package presence

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTestTracker() *Tracker {
	// make a new tracker but clear out the default logging hooks
	t := NewTracker()
	t.onLeaveHooks = nil
	t.onReturnHooks = nil
	return t
}

var (
	dev1 = Device{
		Name: "dev1",
		Addr: "aa:bb:cc:dd:ee:ff",
	}
	dev2 = Device{
		Name: "dev2",
		Addr: "ff:ee:dd:cc:bb:aa",
	}
)

func TestTracker_IsPresent(t *testing.T) {
	ctx := context.Background()

	t.Run("starts in present state", func(t *testing.T) {
		tracker := newTestTracker()
		assert.True(t, tracker.IsPresent())
	})

	t.Run("allows for 2 failures in a device to still be present", func(t *testing.T) {
		tracker := newTestTracker()
		tracker.Set(ctx, dev1, false)
		assert.True(t, tracker.IsPresent())
		tracker.Set(ctx, dev1, false)
		assert.True(t, tracker.IsPresent())
	})

	t.Run("is not present once there are 3 failures", func(t *testing.T) {
		tracker := newTestTracker()
		tracker.Set(ctx, dev1, false)
		tracker.Set(ctx, dev1, false)
		tracker.Set(ctx, dev1, false)
		assert.False(t, tracker.IsPresent())
	})

	t.Run("is present once device reappears", func(t *testing.T) {
		tracker := newTestTracker()
		tracker.Set(ctx, dev1, false)
		tracker.Set(ctx, dev1, false)
		tracker.Set(ctx, dev1, false)
		tracker.Set(ctx, dev1, true)
		assert.True(t, tracker.IsPresent())
	})

	t.Run("requires failures to be consecutive", func(t *testing.T) {
		tracker := newTestTracker()
		tracker.Set(ctx, dev1, false)
		tracker.Set(ctx, dev1, false)
		tracker.Set(ctx, dev1, true)
		tracker.Set(ctx, dev1, false)
		assert.True(t, tracker.IsPresent())

		tracker.Set(ctx, dev1, false)
		tracker.Set(ctx, dev1, false)
		assert.False(t, tracker.IsPresent())
	})

	t.Run("only requires a single device to be missing", func(t *testing.T) {
		tracker := newTestTracker()
		tracker.Set(ctx, dev1, false)
		tracker.Set(ctx, dev2, true)
		tracker.Set(ctx, dev1, false)
		tracker.Set(ctx, dev2, true)
		assert.True(t, tracker.IsPresent())
		tracker.Set(ctx, dev1, false)
		tracker.Set(ctx, dev2, true)
		assert.False(t, tracker.IsPresent())
	})
}

type testHook struct {
	leaveCount  int
	returnCount int
}

func newTestTrackerWithHook() (*Tracker, *testHook) {
	t := newTestTracker()
	h := new(testHook)
	t.OnLeave(h)
	t.OnReturn(h)
	return t, h
}

func (h *testHook) OnLeave(ctx context.Context, t *Tracker) {
	h.leaveCount++
}

func (h *testHook) OnReturn(ctx context.Context, t *Tracker) {
	h.returnCount++
}

func TestTracker_Set(t *testing.T) {
	ctx := context.Background()

	t.Run("trigger leave hook when going from home to away", func(t *testing.T) {
		tracker, hook := newTestTrackerWithHook()
		tracker.Set(ctx, dev1, false)
		tracker.Set(ctx, dev1, false)
		assert.Equal(t, 0, hook.leaveCount)
		tracker.Set(ctx, dev1, false)
		assert.Equal(t, 1, hook.leaveCount)
	})

	t.Run("does not trigger leave hook when already away", func(t *testing.T) {
		tracker, hook := newTestTrackerWithHook()
		tracker.Set(ctx, dev1, false)
		tracker.Set(ctx, dev1, false)
		tracker.Set(ctx, dev1, false)
		assert.Equal(t, 1, hook.leaveCount)
		tracker.Set(ctx, dev1, false)
		assert.Equal(t, 1, hook.leaveCount)
	})

	t.Run("trigger return hook when going from away to home", func(t *testing.T) {
		tracker, hook := newTestTrackerWithHook()
		tracker.Set(ctx, dev1, false)
		tracker.Set(ctx, dev1, false)
		tracker.Set(ctx, dev1, false)
		assert.Equal(t, 0, hook.returnCount)
		tracker.Set(ctx, dev1, true)
		assert.Equal(t, 1, hook.returnCount)
	})

	t.Run("does not trigger return hook when already home", func(t *testing.T) {
		tracker, hook := newTestTrackerWithHook()
		tracker.Set(ctx, dev1, false)
		tracker.Set(ctx, dev1, false)
		tracker.Set(ctx, dev1, false)
		tracker.Set(ctx, dev1, true)
		assert.Equal(t, 1, hook.returnCount)
		tracker.Set(ctx, dev1, true)
		assert.Equal(t, 1, hook.returnCount)
	})
}
