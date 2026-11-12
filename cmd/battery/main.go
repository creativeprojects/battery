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

package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/creativeprojects/battery"
)

func main() {
	batteries, err := battery.GetAll()
	if err, isFatal := err.(battery.ErrFatal); isFatal {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if len(batteries) == 0 {
		fmt.Fprintln(os.Stderr, "No batteries")
		os.Exit(1)
	}
	errs, partialErrs := err.(battery.Errors)
	w := tabwriter.NewWriter(os.Stdout, 2, 2, 2, ' ', tabwriter.AlignRight)
	fmt.Fprintf(w, "Battery\tStatus\tCurrent capacity\tFull capacity\tDesign Capacity\tCharge Rate\tVoltage\t\n")
	for i, batt := range batteries {
		if partialErrs && errs[i] != nil {
			fmt.Fprintf(os.Stderr, "Error getting info for BAT%d: %s\n", i, errs[i])
			continue
		}
		fmt.Fprintf(w, "#%d\t%s\t%.2f mWh\t%.2f mWh\t%.2f mWh\t%.2f mW\t%.2f V\t\n", i, batt.State, batt.Current, batt.Full, batt.Design, batt.ChargeRate, batt.Voltage)
	}
	w.Flush()
}
