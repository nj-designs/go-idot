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
package showimage

import (
	"bytes"
	"fmt"
	"image/png"
	"os"

	"github.com/nj-designs/go-idot/idot"
	"github.com/spf13/cobra"
)

var targetAddr string
var imageFile string

var Cmd = &cobra.Command{
	Use:   "showimage",
	Short: "Shows the supplied .png file on the iDot display",
	Run: func(cmd *cobra.Command, args []string) {
		if err := doShowImage(); err != nil {
			fmt.Printf("error: %v\n", err)
		}
	},
}

func init() {
	Cmd.Flags().StringVar(&targetAddr, "target", "", "Target iDot display MAC address")
	Cmd.MarkFlagRequired("target")

	Cmd.Flags().StringVar(&imageFile, "image-file", "", "Path to a 32x32 .png image file")
	Cmd.MarkFlagRequired("image-file")
}

func validateImage(imageData []byte) error {
	pngImg, err := png.Decode(bytes.NewBuffer(imageData))
	if err != nil {
		return err
	}

	if pngImg.Bounds().Max.X != 32 || pngImg.Bounds().Max.Y != 32 {
		return fmt.Errorf("image is not 32x32")
	}

	return nil
}

func doShowImage() error {
	if len(targetAddr) == 0 {
		return fmt.Errorf("missing --target option")
	}
	if len(imageFile) == 0 {
		return fmt.Errorf("missing --image-file option")
	}

	imageData, err := os.ReadFile(imageFile)
	if err != nil {
		return err
	}
	if err := validateImage(imageData); err != nil {
		return err
	}

	device, err := idot.NewDevice(targetAddr)
	if err != nil {
		return err
	}
	if err := device.Connect(); err != nil {
		return err
	}
	defer device.Disconnect()

	if err := device.SetDrawMode(1); err != nil {
		return err
	}

	if err := device.SendImage(imageData); err != nil {
		return err
	}

	return nil
}
