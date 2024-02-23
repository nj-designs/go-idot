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
	"fmt"

	"tinygo.org/x/bluetooth"
)

const iDotServiceId = uint16(0x00fa)

var iDotServiceUUID = bluetooth.New16BitUUID(iDotServiceId)

const iDotWriteCharacteristicId = uint16(0xfa02)

var iDotWriteCharacteristicUUID = bluetooth.New16BitUUID(iDotWriteCharacteristicId)

const iDotReadCharacteristicId = uint16(0xfa03)

var iDotReadCharacteristicUUID = bluetooth.New16BitUUID(iDotReadCharacteristicId)

var btAdapter = bluetooth.DefaultAdapter

type Device struct {
	scanResult          bluetooth.ScanResult
	btDevice            *bluetooth.Device
	writeCharacteristic bluetooth.DeviceCharacteristic
	writeMTU            int
	readCharacteristic  bluetooth.DeviceCharacteristic
	readMTU             int
}

func NewDevice(targetAddr string) (*Device, error) {
	d := &Device{}

	if err := btAdapter.Enable(); err != nil {
		return d, err
	}

	err := btAdapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
		// println("found device:", result.Address.String(), result.RSSI, result.LocalName())
		if result.Address.String() == targetAddr {
			adapter.StopScan()
			d.scanResult = result
		}
	})
	if err != nil {
		// Scanning didn't start, just return
		return d, err
	}

	return d, nil
}

func (d *Device) Connect() error {

	btd, err := btAdapter.Connect(d.scanResult.Address, bluetooth.ConnectionParams{})
	if err != nil {
		return err
	}

	srvcs, err := btd.DiscoverServices([]bluetooth.UUID{iDotServiceUUID})
	if err != nil {
		return fmt.Errorf("service discover failed")
	}
	if len(srvcs) == 0 {
		return fmt.Errorf("device doesn't support %s service", iDotServiceUUID.String())
	}

	service := srvcs[0]

	if !service.Is16Bit() || service.UUID().Get16Bit() != iDotServiceId {
		return fmt.Errorf("invalid service id")
	}

	chars, err := service.DiscoverCharacteristics([]bluetooth.UUID{iDotWriteCharacteristicUUID, iDotReadCharacteristicUUID})
	if err != nil {
		return err
	}
	if len(chars) != 2 {
		return fmt.Errorf("unexpected number of characteristics. expected 2, got %d", len(chars))
	}

	for _, ch := range chars {
		if !ch.Is16Bit() {
			return fmt.Errorf("invalid char type")
		}
		mtu, err := ch.GetMTU()
		if err != nil {
			return err
		}
		switch ch.Get16Bit() {
		case iDotWriteCharacteristicId:
			d.writeCharacteristic = ch
			d.writeMTU = int(mtu)
		case iDotReadCharacteristicId:
			d.readCharacteristic = ch
			d.readMTU = int(mtu)
		default:
			return fmt.Errorf("invalid characteristic %s", ch.UUID().String())
		}
	}

	d.btDevice = btd

	return nil
}

func (d *Device) Disconnect() {
	d.btDevice.Disconnect()
}

// Write will write the supplied packet to the device
// in up to MTU sized chunks
func (d *Device) Write(packet []byte) error {

	cursor := 0
	remaining := len(packet)
	for remaining > 0 {
		// wl := min(d.writeMTU, remaining)
		wl := min(514, remaining)
		_, err := d.writeCharacteristic.WriteWithoutResponse(packet[cursor : cursor+wl])
		if err != nil {
			return err
		}
		cursor += wl
		remaining -= wl
	}

	return nil
}
