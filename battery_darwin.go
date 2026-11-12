// battery
// Copyright (C) 2016-2017,2023 Karol 'Kenji Takahashi' Woźniak
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the "Software"),
// to deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included
// in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
// DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
// TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
// OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package battery

import (
	"fmt"
	"math"
	"os/exec"

	plist "howett.net/plist"
)

type battery struct {
	BatteryInstalled  bool
	Voltage           int
	CurrentCapacity   int `plist:"AppleRawCurrentCapacity"`
	MaxCapacity       int `plist:"AppleRawMaxCapacity"`
	DesignCapacity    int
	Amperage          int64
	FullyCharged      bool
	IsCharging        bool
	ExternalConnected bool
	// macOS 27 only publishes the raw capacities in this nested dictionary.
	BatteryData batteryData
}

type batteryData struct {
	CurrentCapacity int `plist:"AppleRawCurrentCapacity"`
	MaxCapacity     int `plist:"AppleRawMaxCapacity"`
	DesignCapacity  int
}

// fillCapacities falls back to BatteryData for capacities missing at the top level.
func (b *battery) fillCapacities() {
	if b.CurrentCapacity == 0 {
		b.CurrentCapacity = b.BatteryData.CurrentCapacity
	}
	if b.MaxCapacity == 0 {
		b.MaxCapacity = b.BatteryData.MaxCapacity
	}
	if b.DesignCapacity == 0 {
		b.DesignCapacity = b.BatteryData.DesignCapacity
	}
}

func readBatteries() ([]*battery, error) {
	out, err := exec.Command("ioreg", "-n", "AppleSmartBattery", "-r", "-a").Output()
	if err != nil {
		return nil, err
	}
	return parseBatteries(out)
}

func parseBatteries(raw []byte) ([]*battery, error) {
	var (
		data []*battery
		err  error
	)
	if len(raw) == 0 {
		// No batteries.
		return data, nil
	}

	if _, err = plist.Unmarshal(raw, &data); err != nil {
		return nil, err
	}
	for _, b := range data {
		b.fillCapacities()
	}
	return data, nil
}

func convertBattery(battery *battery) *Battery {
	volts := float64(battery.Voltage) / 1000
	b := &Battery{
		Current:       float64(battery.CurrentCapacity) * volts,
		Full:          float64(battery.MaxCapacity) * volts,
		Design:        float64(battery.DesignCapacity) * volts,
		ChargeRate:    math.Abs(float64(battery.Amperage)) * volts,
		Voltage:       volts,
		DesignVoltage: volts,
	}
	switch {
	case !battery.ExternalConnected:
		b.State.Raw = Discharging
		b.State.specific = "not(ExternalConnected)"
	case battery.IsCharging:
		b.State.Raw = Charging
		b.State.specific = "IsCharging"
	case battery.CurrentCapacity == 0:
		b.State.Raw = Empty
		b.State.specific = "CurrentCapacity is 0"
	case !battery.FullyCharged:
		b.State.Raw = Idle
		b.State.specific = "not(FullyCharged)"
	case battery.FullyCharged:
		b.State.Raw = Full
		b.State.specific = "FullyCharged"
	default:
		b.State.Raw = Undefined
		b.State.specific = fmt.Sprintf("%+v", *battery)
	}
	return b
}

func convertBatteries(source []*battery) []*Battery {
	batteries := make([]*Battery, 0, len(source))
	for _, battery := range source {
		if !battery.BatteryInstalled {
			continue
		}
		batteries = append(batteries, convertBattery(battery))
	}
	return batteries
}

func systemGet(idx int) (*Battery, error) {
	batteries, err := systemGetAll()
	if err != nil {
		return nil, err
	}

	if idx >= len(batteries) {
		return nil, ErrNotFound
	}
	return batteries[idx], nil
}

func systemGetAll() ([]*Battery, error) {
	_batteries, err := readBatteries()
	if err != nil {
		return nil, err
	}

	return convertBatteries(_batteries), nil
}
