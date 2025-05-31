/*
Copyright © 2024 Neil Johnson <nj.designs@protonmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package showclock

import (
	"fmt"
	"time"

	"github.com/nj-designs/go-idot/idot"
	"github.com/spf13/cobra"
)

var clockStyle int
var showDate bool
var show24h bool
var colour string
var timeValue string
var targetAddr string

var Cmd = &cobra.Command{
	Use:   "showclock",
	Short: "Shows and optionally configures the clock of the iDot display",
	Run: func(cmd *cobra.Command, args []string) {
		if err := doSetClock(); err != nil {
			fmt.Printf("error: %v\n", err)
		}
	},
}

func init() {
	Cmd.Flags().StringVar(&targetAddr, "target", "", "Target iDot display MAC address")
	Cmd.MarkFlagRequired("target")
	Cmd.Flags().StringVar(&timeValue, "time", "", "Time value in RFC1123Z format. As per 'date -R'")
	Cmd.Flags().IntVar(&clockStyle, "style", idot.ClockAnimatedHourGlass, "Style of clock. 0:Default 1:Christmas 2:Racing 3:Inverted 4:Hour Glass")
	Cmd.Flags().BoolVar(&showDate, "show-date", true, "Show date as well as time")
	Cmd.Flags().BoolVar(&show24h, "24hour", true, "Show time in 24 hour format")
	Cmd.Flags().StringVar(&colour, "colour", "", "Set RGB colour of clock. Format: R,G,B (0-255)")
}

func doSetClock() error {
	if len(targetAddr) == 0 {
		return fmt.Errorf("missing --target option")
	}

	if clockStyle > idot.ClockAnimatedHourGlass {
		return fmt.Errorf("invalid style")

	}

	var t time.Time
	var err error

	if len(timeValue) > 0 {
		t, err = time.Parse(time.RFC1123Z, timeValue)
		if err != nil {
			return err
		}
	} else {
		t = time.Now()
	}

	device, err := idot.NewDevice(targetAddr)
	if err != nil {
		return err
	}
	if err := device.Connect(); err != nil {
		return err
	}
	defer device.Disconnect()

	if err := device.SetTime(t.Year(), int(t.Month()), t.Day(), int(t.Weekday())+1, t.Hour(),
		t.Minute(), t.Second()); err != nil {
		return err
	}
	customColour, err := idot.ColourFromString(colour)
	if err := device.SetClockMode(clockStyle, showDate, show24h, customColour); err != nil {
		return err
	}

	return nil
}
