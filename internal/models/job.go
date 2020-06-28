package models

import "time"

//Job ...
type Job struct {
	UUID     string
	Download *Download
}

//GetCompletedTime ...
func (j Job) GetCompletedTime() *time.Time {
	return j.Download.Completed
}

//Completed ...
func (j Job) Completed() bool {
	return j.Download.Completed != nil
}

//Running ...
func (j Job) Running() bool {
	return j.Download.Running
}

//Stop ..
func (j Job) Stop() error {
	return j.Download.Stop()
}

//Start ..
func (j Job) Start() error {
	j.Download.Start()
	return nil
}

var jobs = make(map[string]Job)

//AddJob ...
func AddJob(job Job) {
	jobs[job.UUID] = job
}

//AllJobs ...
func AllJobs() map[string]Job {
	return jobs
}
