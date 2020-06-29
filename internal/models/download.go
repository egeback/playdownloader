package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-cmd/cmd"
)

//Download struct containing info regarding a executable command
type Download struct {
	URL             string       `json:"url"`
	Started         *time.Time   `json:"started"`
	Completed       *time.Time   `json:"completed"`
	Running         bool         `json:"running"`
	Done            bool         `json:"done"`
	FindCmd         *cmd.Cmd     `json:"-" swaggerignore:"true"`
	Ticker          *time.Ticker `json:"-" swaggerignore:"true"`
	StdOut          []string     `json:"std_out"`
	StdErr          []string     `json:"std_err"`
	Error           bool         `json:"error"`
	Filename        *string      `json:"filename"`
	AudioETA        *string      `json:"audio_eta"`
	VideoETA        *string      `json:"video_eta"`
	AudioDownloaded bool         `json:"audio_ready"`
	VideoDownloaded bool         `json:"video_ready"`
	AudioStarted    bool         `json:"-" swaggerignore:"true"`
	VideoStarted    bool         `json:"-" swaggerignore:"true"`
	FullStdErr      []string     `json:"-" swaggerignore:"true"`
	Stopped         bool         `json:"stopped"`
}

var ticker *time.Ticker = time.NewTicker(1 * time.Second)

//SvtDLLocation is describing the location of SvtDL
var SvtDLLocation = "/usr/bin/svtplay-dl"

//DefaultDownloadDir specifies default location for downloadedfiles
var DefaultDownloadDir = "/media"

//extractETA from response line
func extractETA(str string) *string {
	i := strings.LastIndex(str, "ETA: ")
	s := str[i+5:]
	return &s
}

//extractFilename from response line
func extractFilename(str string) *string {
	i := strings.Index(str, "into ")
	s := str[i+5:]
	return &s
}

//handleOut parses string for information regarding Outfile and progress
func (d *Download) handleOut(list []string) {
	for _, str := range list {
		//Non progress bar line
		if strings.HasPrefix(str, "INFO:") || strings.HasPrefix(str, "WARNING:") || strings.HasPrefix(str, "DEBUG:") {
			if strings.Index(str, "Outfile:") >= 0 && strings.Index(str, "audio") >= 0 {
				d.AudioStarted = true
			} else if strings.Index(str, "Outfile:") >= 0 && strings.Index(str, "audio") < 0 {
				d.AudioDownloaded = true
				s := "0:00:00"
				d.AudioETA = &s
				d.VideoStarted = true
			} else if strings.Index(str, "INFO: Merge") >= 0 {
				d.AudioDownloaded = true
				d.VideoDownloaded = true
				s := "0:00:00"
				d.AudioETA = &s
				d.VideoETA = &s
			} else if strings.Index(str, "INFO: Muxing ") >= 0 && strings.Index(str, " into ") >= 0 {
				d.Filename = extractFilename(str)
			}
			d.StdErr = append(d.StdErr, str)
		} else if strings.HasPrefix(str, "\r") {
			//Progress bar line. Starts with \r. The go-cmd api gives this as one single string that needs to be splitted
			d.AudioStarted = true
			i := strings.LastIndex(str, "\r")
			statusText := str[i+1:]
			if !strings.HasPrefix(d.StdErr[len(d.StdErr)-1], "[") {
				d.StdErr = append(d.StdErr, statusText)
			} else {
				d.StdErr[len(d.StdErr)-1] = statusText
			}
			if d.VideoStarted {
				d.VideoETA = extractETA(statusText)
			} else {
				d.AudioETA = extractETA(statusText)
			}
		} else {
			//Line does not start with INFO, WARNING, DEBUG or \r
			d.StdErr = append(d.StdErr, str)
		}
	}
}

// Max returns the larger of x or y.
func Max(x int, y int) int {
	if x < y {
		return y
	}
	return x
}

// Start to start download
func (d *Download) Start() {
	d.FindCmd = cmd.NewCmd(SvtDLLocation, d.URL, "--force", "--output", DefaultDownloadDir)
	statusChan := d.FindCmd.Start() // non-blocking
	d.Running = true
	started := time.Now()
	d.Started = &started
	d.Ticker = time.NewTicker(500 * time.Millisecond)
	d.StdErr = make([]string, 0, 10)

	// Print last line of stdout based on ticker interval
	go func() {
		oldNOut := 0
		oldNErr := 0
		for range d.Ticker.C {
			status := d.FindCmd.Status()
			nOut := len(status.Stdout)
			nErr := len(status.Stderr)
			diffOut := nOut - oldNOut

			if len(status.Stderr)-oldNErr > 0 {
				//New lines added from go-cmd
				d.handleOut(status.Stderr[Max(oldNErr, 0):])
			} else if nErr > 0 && strings.HasPrefix(status.Stderr[len(status.Stderr)-1], "\r") {
				//Go-cmd appends progress bar to the current string since now \n is provided from svtplay-dl
				d.handleOut([]string{status.Stderr[nErr-1]})
			} else if nErr > 0 {
				//No new lines added between ticks
			}

			//svtplay-dl normaly prints to SdErr but for debuging purposes we store the stdout
			if diffOut > 0 {
				str := status.Stdout[nOut]
				fmt.Println(str)
			}

			//Update nOut and nErr
			oldNOut = nOut
			oldNErr = nErr
			d.StdOut = status.Stdout

			//Break if status is go-cmd inidcates that the process has closed
			if status.Complete {
				break
			}
		}
	}()

	// Stop command after 1 hour
	go func() {
		<-time.After(1 * time.Hour)
		d.FindCmd.Stop()
		d.Error = true
		t := time.Now()
		d.Completed = &t
	}()

	// Check if command is done before the ticker got called
	select {
	case <-statusChan:
		t := time.Now()
		d.Completed = &t
		d.Done = true
		d.FullStdErr = d.FindCmd.Status().Stderr
	default:
		// no, still running
		go func() {
			// Block waiting for command to exit, be stopped, or be killed
			<-statusChan
			d.Done = true
			d.Running = false
			t := time.Now()
			d.Completed = &t
			d.FullStdErr = d.FindCmd.Status().Stderr
		}()
	}
}

//CreateDownload struct and return it
func CreateDownload(url string) Download {
	return Download{URL: url, Running: false, Error: false, Done: false, AudioDownloaded: false, VideoDownloaded: false, Stopped: false}
}

//Stop download
func (d *Download) Stop() error {
	err := d.FindCmd.Stop()
	if err == nil {
		t := time.Now()
		d.Completed = &t
		d.Stopped = true
	}
	return err
}
