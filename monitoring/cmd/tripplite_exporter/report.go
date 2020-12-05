package main

import (
	"fmt"

	"github.com/zserge/hid"
)

const (
	ReportConfigVoltage   = 1
	ReportConfigFrequency = 2
	ReportConfigPower     = 3
	ReportInputVoltage    = 24
	ReportInputFrequency  = 25
	ReportOutputVoltage   = 27
	ReportStatus          = 35
	ReportHealth          = 52
	ReportTimeToEmpty     = 53
	ReportOutputPower     = 71
)

func ReadInt8(device hid.Device, feature int) (int8, error) {
	data, err := readFeature(device, feature)
	if err != nil {
		return 0, err
	}
	return int8(data[0]), nil
}

func ReadInt16(device hid.Device, feature int) (int16, error) {
	data, err := readFeature(device, feature)
	if err != nil {
		return 0, err
	}
	return (int16(data[1]) << 8) + int16(data[0]), nil
}

func ReadFloat(device hid.Device, feature int) (float64, error) {
	val, err := ReadInt16(device, feature)
	if err != nil {
		return 0, err
	}
	return float64(val) / 10.0, nil
}

func ReadBitSet(device hid.Device, feature int) ([]bool, error) {
	val, err := ReadInt8(device, feature)
	if err != nil {
		return nil, err
	}

	var flags [8]bool
	for i := 0; i < 8; i++ {
		flags[i] = val&(1<<i) != 0
	}
	return flags[:], nil
}

func readFeature(device hid.Device, feature int) ([]byte, error) {
	data, err := device.GetReport(feature)
	if err != nil {
		return nil, err
	}
	if data[0] != byte(feature) {
		return nil, fmt.Errorf("unexpected data in feature report")
	}
	return data[1:], nil
}