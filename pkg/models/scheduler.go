package models

import (
	"log"
	"time"
)

//DefaultNrWorkers ..
var DefaultNrWorkers = 2

//DefaultCompletedJobRemoval ..
var DefaultCompletedJobRemoval = time.Duration(-15) * time.Minute

//Scheduler ...
type Scheduler struct {
	Workers             int
	jobs                []Job
	completedJobRemoval time.Duration
	ticker              *time.Ticker
}

//NewScheduler ...
func NewScheduler() *Scheduler {
	s := &Scheduler{
		Workers:             DefaultNrWorkers,
		jobs:                make([]Job, 0, 0),
		completedJobRemoval: DefaultCompletedJobRemoval,
	}
	s.Start()
	return s
}

//GetJobs ...
func (s *Scheduler) GetJobs() []Job {
	return s.jobs
}

//GetJobByUUID ...
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

//GetNrRunningJobs ..
func (s *Scheduler) getNrRunningJobs() int {
	nrRunning := 0
	for _, j := range s.jobs {
		if j.Running() {
			nrRunning++
		}
	}
	return nrRunning
}

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

//Stop ...
func (s *Scheduler) Stop() {
	s.ticker.Stop()
	for _, j := range s.jobs {
		if err := j.Stop(); err != nil {
			log.Panic(err)
		}
	}
}

//Start ...
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
				//then := time.Now().Add(s.completedJobRemoval)
				//if j.GetCompletedTime().Before(then) {
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
