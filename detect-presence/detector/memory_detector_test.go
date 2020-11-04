package detector

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoryDetector_IsHealthy(t *testing.T) {
	ctx := context.Background()

	t.Run("starts healthy", func(t *testing.T) {
		d := NewMemoryDetector("aa:bb:cc:dd:ee:ff")
		healthy, err := d.IsHealthy(ctx)
		assert.NoError(t, err)
		assert.True(t, healthy)
	})

	t.Run("can report as unhealthy", func(t *testing.T) {
		d := NewMemoryDetector("aa:bb:cc:dd:ee:ff")
		d.SetBluetoothHealth(false)
		healthy, err := d.IsHealthy(ctx)
		assert.NoError(t, err)
		assert.False(t, healthy)
	})

	t.Run("can report an error", func(t *testing.T) {
		d := NewMemoryDetector("aa:bb:cc:dd:ee:ff")
		d.SetBluetoothError(fmt.Errorf("error checking bluetooth"))
		_, err := d.IsHealthy(ctx)
		assert.Error(t, err)
	})
}

func TestMemoryDetector_DetectDevice(t *testing.T) {
	ctx := context.Background()

	t.Run("starts present", func(t *testing.T) {
		d := NewMemoryDetector("aa:bb:cc:dd:ee:ff")
		present, err := d.DetectDevice(ctx, "aa:bb:cc:dd:ee:ff")
		assert.NoError(t, err)
		assert.True(t, present)
	})

	t.Run("reports error for unknown devices", func(t *testing.T) {
		d := NewMemoryDetector("aa:bb:cc:dd:ee:ff")
		_, err := d.DetectDevice(ctx, "aa:aa:aa:aa:aa:aa")
		assert.Error(t, err)
	})

	t.Run("can report not present", func(t *testing.T) {
		d := NewMemoryDetector("aa:bb:cc:dd:ee:ff")
		d.SetDevicePresence("aa:bb:cc:dd:ee:ff", false)
		present, err := d.DetectDevice(ctx, "aa:bb:cc:dd:ee:ff")
		assert.NoError(t, err)
		assert.False(t, present)
	})

	t.Run("can report an error", func(t *testing.T) {
		d := NewMemoryDetector("aa:bb:cc:dd:ee:ff")
		d.SetDeviceError("aa:bb:cc:dd:ee:ff", fmt.Errorf("error checking device"))
		_, err := d.DetectDevice(ctx, "aa:bb:cc:dd:ee:ff")
		assert.Error(t, err)
	})

	t.Run("supports multiple devices", func(t *testing.T) {
		d := NewMemoryDetector("aa:bb:cc:dd:ee:ff", "ff:ee:dd:cc:bb:aa")

		d.SetDevicePresence("aa:bb:cc:dd:ee:ff", false)
		present, err := d.DetectDevice(ctx, "aa:bb:cc:dd:ee:ff")
		assert.NoError(t, err)
		assert.False(t, present)
		present, err = d.DetectDevice(ctx, "ff:ee:dd:cc:bb:aa")
		assert.NoError(t, err)
		assert.True(t, present)

		d.SetDeviceError("ff:ee:dd:cc:bb:aa", fmt.Errorf("error checking device"))
		present, err = d.DetectDevice(ctx, "aa:bb:cc:dd:ee:ff")
		assert.NoError(t, err)
		assert.False(t, present)
		_, err = d.DetectDevice(ctx, "ff:ee:dd:cc:bb:aa")
		assert.Error(t, err)
	})
}
