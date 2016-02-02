package config

import (
	"encoding/json"
	"time"

	log "github.com/Sirupsen/logrus"
)

// Configuration stores list of a Jobs to run
type Configuration struct {
	Jobs []Job
}

// Job is an item in Configuration
type Job struct {
	Name             string
	Command          string
	Args             []string
	AllowedInstances int
	Timeout          int
	Runners          int
	RunsEvery        string
	duration         time.Duration
	Instances        int
}

func (job *Job) Interval() time.Duration {
	if job.duration == 0 {
		duration, err := time.ParseDuration(job.RunsEvery)

		if err != nil {
			panic(err)
		}

		job.duration = duration
	}

	return job.duration
}

func (job *Job) IncrementInstances() {
	job.Instances++
}

// @TODO add checking for below 0
func (job *Job) DecrementInstances() {
	job.Instances--
}

// @TODO add checking for below 0
func (job *Job) CanRun() bool {
	log.Debug("", job.Instances, job.Runners)
	return (job.Instances < job.Runners)
}

// Parse parses json config into Configuration object
func Parse(jsonConfig []byte) (*Configuration, error) {
	c := &Configuration{}
	err := json.Unmarshal(jsonConfig, &c)
	if err != nil {
		return nil, err
	}

	for _, job := range c.Jobs {
		log.Debug("", job.RunsEvery, job.Instances, job.Runners)
	}

	return c, nil
}
