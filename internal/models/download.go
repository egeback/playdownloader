package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-cmd/cmd"
)

var ticker *time.Ticker = time.NewTicker(1 * time.Second)

//SvtDLLocation is describing the location of SvtDL
var SvtDLLocation = "/usr/bin/svtplay-dl"

//DefaultDownloadDir specifies default location for downloadedfiles
var DefaultDownloadDir = "/media"

//Download struct containing info regarding a executable command
type Download struct {
	URL             string            `json:"url"`
	Started         *time.Time        `json:"started"`
	Completed       *time.Time        `json:"completed"`
	Running         bool              `json:"running"`
	Done            bool              `json:"done"`
	FindCmd         *cmd.Cmd          `json:"-"`
	StatusChan      <-chan cmd.Status `json:"-"`
	Ticker          *time.Ticker      `json:"-"`
	StdOut          []string          `json:"std_out"`
	StdErr          []string          `json:"std_err"`
	Error           bool              `json:"error"`
	Filename        *string           `json:"filename"`
	AudioETA        *string           `json:"audio_eta"`
	VideoETA        *string           `json:"video_eta"`
	AudioDownloaded bool              `json:"audio_ready"`
	VideoDownloaded bool              `json:"video_ready"`
	AudioStarted    bool              `json:"-"`
	VideoStarted    bool              `json:"-"`
	FullStdErr      []string          `json:"-"`
}

func extractETA(str string) *string {
	i := strings.LastIndex(str, "ETA: ")
	s := str[i+5:]
	return &s
}

func extractFilename(str string) *string {
	i := strings.Index(str, "into ")
	s := str[i+5:]
	return &s
}

func (d *Download) handleOut(list []string) {
	for _, str := range list {
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
	d.StatusChan = d.FindCmd.Start() // non-blocking
	d.Running = true
	started := time.Now()
	d.Started = &started
	d.Ticker = time.NewTicker(500 * time.Millisecond)
	d.StdErr = make([]string, 0, 10)

	// Print last line of stdout every 2s
	go func() {
		oldNOut := 0
		oldNErr := 0
		for range d.Ticker.C {
			status := d.FindCmd.Status()
			nOut := len(status.Stdout)
			nErr := len(status.Stderr)
			diffOut := nOut - oldNOut

			if len(status.Stderr)-oldNErr > 0 {
				d.handleOut(status.Stderr[Max(oldNErr, 0):])
			} else if nErr > 0 && strings.HasPrefix(status.Stderr[len(status.Stderr)-1], "\r") {
				d.handleOut([]string{status.Stderr[nErr-1]})
			} else if nErr > 0 {
				if len(status.Stderr)-oldNErr > 0 {
					fmt.Println("It happend", len(status.Stderr)-oldNErr)
				}
			}

			if diffOut > 0 {
				str := status.Stdout[nOut]
				fmt.Println(str)
			}

			oldNOut = nOut
			oldNErr = nErr
			d.StdOut = status.Stdout

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

	// Check if command is done
	select {
	case <-d.StatusChan:
		t := time.Now()
		d.Completed = &t
		d.Done = true
		d.FullStdErr = d.FindCmd.Status().Stderr
	default:
		// no, still running
		go func() {
			// Block waiting for command to exit, be stopped, or be killed
			<-d.StatusChan
			d.Done = true
			d.Running = false
			t := time.Now()
			d.Completed = &t
			d.FullStdErr = d.FindCmd.Status().Stderr
		}()
	}
}

var downloads []Download

//DownloadsAll to retrive all downloders
func DownloadsAll() []Download {
	return downloads
}

//AddDownload ..
func AddDownload(url string) Download {
	if downloads == nil {
		downloads = make([]Download, 0, 10)
	}
	download := Download{URL: url, Running: false, Error: false, Done: false, AudioDownloaded: false, VideoDownloaded: false}
	downloads = append(downloads, download)
	return download
}

//Stop ...
func (d Download) Stop() error {
	return d.FindCmd.Stop()
}
