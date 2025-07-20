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
package startserver

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/nj-designs/go-idot/idot"
	"github.com/spf13/cobra"
)

type iDotService struct {
	device *idot.Device
}

var serverPort uint
var targetAddr string

const apiBase = "/api/v1"

var Cmd = &cobra.Command{
	Use:   "startserver",
	Short: "Start a simple rest API server",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runServer(); err != nil {
			fmt.Printf("Failed: %v\n", err)
		}
	},
}

func init() {

	Cmd.Flags().StringVar(&targetAddr, "target", "", "Target iDot display MAC address")
	Cmd.MarkFlagRequired("target")

	Cmd.Flags().UintVar(&serverPort, "port", 8080, "Port to listen on")
}

func runServer() error {

	device, err := idot.NewDevice(targetAddr)
	if err != nil {
		return err
	}
	fmt.Printf("Connecting to %s\n", targetAddr)
	if err := device.Connect(); err != nil {
		return err
	}
	defer device.Disconnect()
	fmt.Println("Connected")

	ids := &iDotService{device: device}

	mux := http.NewServeMux()
	mux.HandleFunc(fmt.Sprintf("POST %s", formFullUrl("/showclock/")), ids.handleShowClock)
	mux.HandleFunc(fmt.Sprintf("POST %s", formFullUrl("/showimage/")), ids.handleShowImage)

	srv := &http.Server{Addr: fmt.Sprintf(":%d", serverPort), Handler: mux}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs
		fmt.Println("\nStart shutdown")
		if err := srv.Shutdown(context.Background()); err != nil {
			fmt.Printf("Shutdown returned: %v\n", err)
		}
		close(idleConnsClosed)
	}()

	fmt.Printf("Listing at %s\n", srv.Addr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	<-idleConnsClosed

	return nil
}

func formFullUrl(endPoint string) string {
	return path.Join(apiBase, endPoint)
}

type setClockValues struct {
	Time     string `json:"time,omitempty"`
	Style    int    `json:"style,omitempty"`
	ShowDate bool   `json:"showdate,omitempty"`
	Show24h  bool   `json:"show24h,omitempty"`
	Colour  string   `json:"colour,omitempty"`
}

func (ids *iDotService) handleShowClock(w http.ResponseWriter, req *http.Request) {
	cv := &setClockValues{}
	if req.ContentLength > 0 {
		if err := json.NewDecoder(req.Body).Decode(cv); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	// fmt.Println(cv)
	var t time.Time
	var err error

	if len(cv.Time) > 0 {
		t, err = time.Parse(time.RFC1123Z, cv.Time)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		t = time.Now()
	}
	if err := ids.device.SetTime(t.Year(), int(t.Month()), t.Day(), int(t.Weekday())+1, t.Hour(),
		t.Minute(), t.Second()); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	customColour, err := idot.ColourFromString(cv.Colour)
	if err := ids.device.SetClockMode(cv.Style, cv.ShowDate, cv.Show24h, customColour); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (ids *iDotService) handleShowImage(w http.ResponseWriter, req *http.Request) {
	req.ParseMultipartForm(1024 * 1024)
	file, handler, err := req.FormFile("imgfile")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)
	fileData, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := ids.device.SetDrawMode(1); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := ids.device.SendImage(fileData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
