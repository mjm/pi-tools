package detector

import (
	"context"
	"fmt"
)

type MemoryDetector struct {
	bt      checkStatus
	devices map[string]checkStatus
}

type checkStatus struct {
	present bool
	err     error
}

func NewMemoryDetector(addrs ...string) *MemoryDetector {
	devices := map[string]checkStatus{}
	for _, addr := range addrs {
		devices[addr] = checkStatus{present: true}
	}
	return &MemoryDetector{
		bt:      checkStatus{present: true},
		devices: devices,
	}
}

func (d *MemoryDetector) SetBluetoothError(err error) {
	d.bt.err = err
}

func (d *MemoryDetector) SetBluetoothHealth(healthy bool) {
	d.bt = checkStatus{present: healthy}
}

func (d *MemoryDetector) SetDeviceError(addr string, err error) {
	d.devices[addr] = checkStatus{err: err}
}

func (d *MemoryDetector) SetDevicePresence(addr string, present bool) {
	d.devices[addr] = checkStatus{present: present}
}

func (d *MemoryDetector) IsHealthy(ctx context.Context) (bool, error) {
	if d.bt.err != nil {
		return false, d.bt.err
	}

	return d.bt.present, nil
}

func (d *MemoryDetector) DetectDevice(ctx context.Context, addr string) (bool, error) {
	dev, ok := d.devices[addr]
	if !ok {
		return false, fmt.Errorf("no such device %q", addr)
	}

	if dev.err != nil {
		return false, dev.err
	}

	return dev.present, nil
}
