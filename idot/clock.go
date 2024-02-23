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
package idot

const (
	ClockDefault           = iota
	ClockChristmas         = iota
	ClockRacing            = iota
	ClockInverted          = iota
	ClockAnimatedHourGlass = iota
)

func (d *Device) SetClockMode(style int, visibleDate bool, hour24 bool, colour Colour) error {
	var sb uint8 = uint8(style)
	if visibleDate {
		sb |= 128
	}
	if hour24 {
		sb |= 64
	}
	return d.Write([]byte{8, 0, 6, 1, sb, colour.R, colour.G, colour.B})
}

func (d *Device) SetTime(year int, month int, day int, weekDay int, hour int, minute int, second int) error {

	return d.Write([]byte{11, 0, 1, 128, byte(year), byte(month), byte(day), byte(weekDay), byte(hour), byte(minute), byte(second)})
}
