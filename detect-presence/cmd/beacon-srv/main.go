package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os/exec"

	"github.com/google/uuid"

	"github.com/mjm/pi-tools/pkg/signal"
)

var (
	device           = flag.String("device", "hci0", "Bluetooth device name to use to advertise")
	proximityUUIDStr = flag.String("proximity-uuid", "", "Proximity UUID to advertise for iBeacon")
	major            = flag.Uint("major", 0, "Major value to advertise for iBeacon")
	minor            = flag.Uint("minor", 0, "Minor value to advertise for iBeacon")
)

func main() {
	flag.Parse()

	if *proximityUUIDStr == "" {
		log.Panicf("no proximity UUID provided")
	}

	proximityUUID, err := uuid.Parse(*proximityUUIDStr)
	if err != nil {
		log.Panicf("failed to parse proximity UUID: %v", err)
	}

	var idBytes []string
	for _, b := range proximityUUID {
		idBytes = append(idBytes, hex.EncodeToString([]byte{b}))
	}
	idBytes = append(idBytes, fmt.Sprintf("%02x", *major>>8), fmt.Sprintf("%02x", *major&0xFF))
	idBytes = append(idBytes, fmt.Sprintf("%02x", *minor>>8), fmt.Sprintf("%02x", *minor&0xFF))

	log.Printf("enabling bluetooth advertising mode")
	if err := exec.Command("hciconfig", *device, "leadv", "3").Run(); err != nil {
		// Exiting with status 1 is actually okay, it may mean we've already set this
		if _, ok := err.(*exec.ExitError); !ok {
			log.Panicf("failed to enable bluetooth advertising mode: %v", err)
		}
	}

	log.Printf("disabling bluetooth scanning")
	out, err := exec.Command("hciconfig", *device, "noscan").CombinedOutput()
	if err != nil {
		log.Panicf("failed to disable bluetooth scanning: %v\n%s", err, out)
	}

	log.Printf("setting advertising data")
	args := []string{
		"-i", *device, "cmd", "0x08", "0x0008",
		"1E", "02", "01", "06", "1A", "FF", "4C", "00", "02", "16", "15",
	}
	args = append(args, idBytes...)
	args = append(args, "C8", "00")
	out, err = exec.Command("hcitool", args...).CombinedOutput()
	if err != nil {
		log.Panicf("failed to set advertising data: %v\n%s", err, out)
	}

	signal.Wait()
}
