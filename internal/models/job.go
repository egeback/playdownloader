package models

import "time"

//Job struct to including Download object
type Job struct {
	UUID     string    `json:"uuid" format:"uuid"`
	Download *Download `json:"download_info"`
}

//GetCompletedTime return Download.Completed
func (j Job) GetCompletedTime() *time.Time {
	return j.Download.Completed
}

//Completed evaluate if Download.Completed is set
func (j Job) Completed() bool {
	return j.Download.Completed != nil
}

//Running provides Download.Running
func (j Job) Running() bool {
	return j.Download.Running
}

//Stop stop download
func (j Job) Stop() error {
	return j.Download.Stop()
}

//Start job
func (j Job) Start() error {
	j.Download.Start()
	return nil
}
