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

import (
	"errors"
    "strconv"
    "strings"
)

type Colour struct {
	R, G, B uint8
}

var Red = Colour{255, 0, 0}
var Green = Colour{0, 255, 0}
var Blue = Colour{0, 0, 255}



func ColourFromString(rgb string) (Colour, error) {
    parts := strings.Split(rgb, ",")
    if len(parts) != 3 {
        return Colour{}, ErrInvalidRGB
    }
    r, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
    g, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
    b, err3 := strconv.Atoi(strings.TrimSpace(parts[2]))
    if err1 != nil || err2 != nil || err3 != nil {
        return Colour{}, ErrInvalidRGB
    }
    return Colour{uint8(r), uint8(g), uint8(b)}, nil
}

var ErrInvalidRGB = errors.New("Invalid RGB value. Please use the format 'R, G, B'.")