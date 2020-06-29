package models

import (
	"log"
	"time"
)

//DefaultNrWorkers determines how many downloads will run at the same time
var DefaultNrWorkers = 2

//DefaultCompletedJobRemoval determines how long a until a job is deleted
var DefaultCompletedJobRemoval = time.Duration(-15) * time.Minute

//Scheduler struct to handle scheduler data
type Scheduler struct {
	Workers             int
	jobs                []Job
	completedJobRemoval time.Duration
	ticker              *time.Ticker
}

//NewScheduler creates new Scheduler
func NewScheduler() *Scheduler {
	s := &Scheduler{
		Workers:             DefaultNrWorkers,
		jobs:                make([]Job, 0, 0),
		completedJobRemoval: DefaultCompletedJobRemoval,
	}
	s.Start()
	return s
}

//GetJobs return all jobs
func (s *Scheduler) GetJobs() []Job {
	return s.jobs
}

//GetJobByUUID return job by string uuid
func (s *Scheduler) GetJobByUUID(uuid string) *Job {
	for _, j := range s.jobs {
		if j.UUID == uuid {
			return &j
		}
	}
	return nil
}

// AddJob starts goroutine which constantly calls provided job with interval delay
func (s *Scheduler) AddJob(job Job) {
	s.jobs = append(s.jobs, job)
}

// AllJobs returns job from scheduler
func (s *Scheduler) AllJobs(uuid string) []Job {
	return s.jobs
}

//GetNrRunningJobs returns nummer of running jobs by iterating and checking Running parameter
func (s *Scheduler) getNrRunningJobs() int {
	nrRunning := 0
	for _, j := range s.jobs {
		if j.Running() {
			nrRunning++
		}
	}
	return nrRunning
}

//getStartableJobs returns a list of jobs that can be started based on nr parameter and number of available non started jobs
func (s *Scheduler) getStartableJobs(nr int) []Job {
	n := 0
	jobs := make([]Job, 0, nr)
	for _, j := range s.jobs {
		if !j.Running() && !j.Completed() {
			jobs = append(jobs, j)
			n++
			if n >= nr {
				return jobs
			}
		}
	}
	return jobs
}

//Stop stop scheduler and stopp all running jobs
func (s *Scheduler) Stop() {
	s.ticker.Stop()
	for _, j := range s.jobs {
		if err := j.Stop(); err != nil {
			log.Println("Error closing job: ", err)
		}
	}
}

//Start scheduler and ticker
func (s *Scheduler) Start() {
	s.ticker = time.NewTicker(500 * time.Second)

	go func() {
		for range ticker.C {
			//1. Check nr running jobs
			nr := s.Workers - s.getNrRunningJobs()

			//2. Start jobs
			if nr > 0 {
				for _, j := range s.getStartableJobs(nr) {
					j.Start()
				}
			}

			//3. Remove old jobs
			for i, j := range s.jobs {
				if j.GetCompletedTime() == nil {
					continue
				}
				diff := time.Now().Sub(*j.GetCompletedTime())
				if diff.Minutes() >= 15 {
					s.jobs = append(s.jobs[:i], s.jobs[i+1:]...)
				}
			}
		}
	}()
}
