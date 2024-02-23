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
package btscan

import (
	"fmt"
	"math"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"tinygo.org/x/bluetooth"
)

var maxScanTime uint32
var verbose bool

var Cmd = &cobra.Command{
	Use:   "btscan",
	Short: "Displays a list of bluetooth devices that can be seen by the local adapter",
	Run: func(cmd *cobra.Command, args []string) {
		if err := doBTScan(); err != nil {
			fmt.Printf("Failed: %v\n", err)
		}
	},
}

func init() {
	Cmd.Flags().Uint32Var(&maxScanTime, "scan-time", 0, "Max number of seconds to perform scan. 0 means infinite")
	Cmd.Flags().BoolVar(&verbose, "verbose", false, "Verbose output during scan")
}

func doBTScan() error {

	if maxScanTime == 0 {
		maxScanTime = math.MaxUint32
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	adapter := bluetooth.DefaultAdapter

	if err := adapter.Enable(); err != nil {
		return err
	}

	scanTimer := time.NewTimer(time.Second * time.Duration(maxScanTime))
	go func() {
		select {
		case <-sigs:
			scanTimer.Stop()
		case <-scanTimer.C:
		}
		adapter.StopScan()
	}()

	scanResults := make(map[string]bluetooth.ScanResult)

	if maxScanTime == math.MaxUint32 {
		fmt.Println("Scanning forever [CTRL+C to stop]")
	} else {
		fmt.Printf("Scanning for %d second(s)\n", maxScanTime)
	}

	err := adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
		addr := result.Address.String()
		if _, prs := scanResults[addr]; !prs {
			if verbose {
				fmt.Printf("Found Device at %s\n", addr)
			}
			scanResults[addr] = result
		}
	})
	if err != nil {
		return err
	}

	fmt.Println("Scan results")
	for _, sr := range scanResults {
		fmt.Printf("Address:%s  RSSI:%3d  Name:%s\n", sr.Address.String(), sr.RSSI, sr.LocalName())
	}

	return nil
}
