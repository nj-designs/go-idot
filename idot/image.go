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
	"bytes"
	"encoding/binary"
)

// SetDrawMode sends set draw mode to display
func (d *Device) SetDrawMode(mode int) error {
	return d.Write([]byte{5, 0, 4, 1, uint8(mode)})
}

// SendImage sends an image to the display. Only makes sense after a call to SetDrawMode(1)
func (d *Device) SendImage(imageData []byte) error {

	// Based on create_payloads in core/idotmatrix/image.py
	chunks := chunkBuffer(imageData, 4096)
	cib := new(bytes.Buffer)
	idk := len(imageData) + len(chunks)
	for ci, ch := range chunks {
		binary.Write(cib, binary.LittleEndian, uint16(idk)) //struct.pack("h", idk)
		binary.Write(cib, binary.LittleEndian, uint8(0))
		binary.Write(cib, binary.LittleEndian, uint8(0))
		if ci > 0 {
			binary.Write(cib, binary.LittleEndian, uint8(2))
		} else {
			binary.Write(cib, binary.LittleEndian, uint8(0))
		}
		binary.Write(cib, binary.LittleEndian, int32(len(imageData))) // struct.pack("i", len(png_data))
		binary.Write(cib, binary.LittleEndian, ch)
	}

	return d.Write(cib.Bytes())
}

// chunkBuffer chunks the supplied data buffer to chunkSize slices
func chunkBuffer(data []byte, chunkSize int) [][]byte {
	chunks := make([][]byte, 0)

	cursor := 0
	remaining := len(data)
	for remaining > 0 {
		wl := min(chunkSize, remaining)
		chunks = append(chunks, data[cursor:cursor+wl])
		cursor += wl
		remaining -= wl
	}
	return chunks
}
