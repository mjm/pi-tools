package presence

// Device represents a single physical device whose presence can be checked for.
type Device struct {
	// Name is an arbitrary label given to identify the device.
	Name string
	// Addr is the Bluetooth address of the device.
	Addr string
	// Canary is true if the device is intended to function as a canary. A canary device is expected to always be
	// present. If a canary device is determined to not be present, then a metric will be set to trigger an alert, and
	// no other devices will be checked for presence until the canary device reports as present.
	Canary bool
}
