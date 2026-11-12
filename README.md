battery [![Build Status](https://github.com/creativeprojects/battery/actions/workflows/tests.yml/badge.svg?branch=master)](https://github.com/creativeprojects/battery/actions/workflows/tests.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/creativeprojects/battery)](https://goreportcard.com/report/github.com/creativeprojects/battery) [![GoDoc](https://pkg.go.dev/badge/github.com/creativeprojects/battery)](https://pkg.go.dev/github.com/creativeprojects/battery)
=======

Cross-platform, normalized battery information library.

Gives access to a system independent, typed battery state, capacity, charge and voltage values recalculated as necessary to be returned in mW, mWh or V units.

Currently supported systems:

* Linux
* macOS (Intel and Apple Silicon, up to macOS 27)
* Windows
* FreeBSD
* DragonFlyBSD
* NetBSD
* OpenBSD
* Solaris

Requires Go 1.26 or later.

Installation
------------

```bash
$ go get github.com/creativeprojects/battery
```

Code Example
------------

```go
package main

import (
	"fmt"

	"github.com/creativeprojects/battery"
)

func main() {
	batteries, err := battery.GetAll()
	if err != nil {
		fmt.Println("Could not get battery info!")
		return
	}
	for i, bat := range batteries {
		fmt.Printf("Bat%d: ", i)
		fmt.Printf("state: %s, ", bat.State)
		fmt.Printf("current capacity: %f mWh, ", bat.Current)
		fmt.Printf("last full capacity: %f mWh, ", bat.Full)
		fmt.Printf("design capacity: %f mWh, ", bat.Design)
		fmt.Printf("charge rate: %f mW, ", bat.ChargeRate)
		fmt.Printf("voltage: %f V, ", bat.Voltage)
		fmt.Printf("design voltage: %f V\n", bat.DesignVoltage)
	}
}
```

`GetAll()` returns an empty slice when the system has no battery attached. Use `Get(idx)` to read a single battery.

Battery state
-------------

`Battery.State.Raw` is one of the platform agnostic states: `Unknown`, `Empty`, `Full`, `Charging`, `Discharging`, `Idle` (plugged in but holding its charge, e.g. capped at 80%) or `Undefined` (the system reported a state the library cannot map).

`Battery.State.Explain()` returns the platform specific value the state was derived from.

Error handling
--------------

* `ErrFatal`: no meaningful data could be retrieved.
* `ErrPartial`: some fields could not be read; it holds one error per `Battery` field, the other fields are valid.
* `Errors`: returned by `GetAll()` only, the i-th entry (`nil`, `ErrPartial` or `ErrFatal`) matches the i-th battery.

`Get()` returns either `ErrFatal` or `ErrPartial`, `GetAll()` returns either `ErrFatal` or `Errors`.

CLI
---

There is also a little utility which prints a summary of all the batteries.

*Installation*

```bash
$ go install github.com/creativeprojects/battery/cmd/battery@latest
```

*Usage*

```bash
$ battery
  Battery       Status  Current capacity  Full capacity  Design Capacity  Charge Rate  Voltage
       #0  Discharging      57466.37 mWh   62945.84 mWh     75141.68 mWh   4786.80 mW  12.37 V
```
