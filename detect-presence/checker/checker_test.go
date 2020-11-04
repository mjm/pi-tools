package checker

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"

	"github.com/mjm/pi-tools/detect-presence/detector"
	"github.com/mjm/pi-tools/detect-presence/presence"
)

var (
	dev1 = presence.Device{
		Name: "dev1",
		Addr: "aa:bb:cc:dd:ee:ff",
	}
	dev2 = presence.Device{
		Name: "dev2",
		Addr: "ff:ee:dd:cc:bb:aa",
	}
)

func TestChecker_Run(t *testing.T) {
	clock := clockwork.NewFakeClock()

	t.Run("tracks health of bluetooth device", func(t *testing.T) {
		d := detector.NewMemoryDetector(dev1.Addr)
		tracker := presence.NewTracker()
		c := &Checker{
			Tracker:  tracker,
			Detector: d,
			Interval: 30 * time.Second,
			Devices:  []presence.Device{dev1},
			clock:    clock,
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		tickCh := make(chan struct{})
		go c.Run(ctx, tickCh)

		tick := func() {
			clock.Advance(30 * time.Second)
			<-tickCh
		}

		<-tickCh
		assert.True(t, c.isHealthy)

		d.SetBluetoothHealth(false)
		assert.True(t, c.isHealthy)
		tick()
		assert.False(t, c.isHealthy)

		d.SetBluetoothHealth(true)
		assert.False(t, c.isHealthy)
		tick()
		assert.True(t, c.isHealthy)

		d.SetBluetoothError(fmt.Errorf("error checking bluetooth"))
		assert.True(t, c.isHealthy)
		tick()
		assert.False(t, c.isHealthy)

		d.SetBluetoothHealth(true)
		assert.False(t, c.isHealthy)
		tick()
		assert.True(t, c.isHealthy)
	})

	t.Run("updates device status in the tracker", func(t *testing.T) {
		d := detector.NewMemoryDetector(dev1.Addr, dev2.Addr)
		tracker := presence.NewTracker()
		c := &Checker{
			Tracker:  tracker,
			Detector: d,
			Interval: 30 * time.Second,
			Devices:  []presence.Device{dev1, dev2},
			clock:    clock,
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		tickCh := make(chan struct{})
		go c.Run(ctx, tickCh)

		tick := func() {
			clock.Advance(30 * time.Second)
			<-tickCh
		}

		<-tickCh
		assert.True(t, tracker.IsPresent())

		d.SetDevicePresence(dev1.Addr, false)
		tick()
		tick()
		assert.True(t, tracker.IsPresent())
		tick()
		assert.False(t, tracker.IsPresent())

		d.SetDevicePresence(dev2.Addr, false)
		tick()
		assert.False(t, tracker.IsPresent())
		tick() // ensure this one ticks enough times to surpass the allowed failure count

		d.SetDevicePresence(dev1.Addr, true)
		tick()
		assert.False(t, tracker.IsPresent())

		d.SetDevicePresence(dev2.Addr, true)
		assert.False(t, tracker.IsPresent())
		tick()
		assert.True(t, tracker.IsPresent())
	})
}
