package gron

import (
	"errors"
	"log"
	"strings"
	"sync"
	"time"
	"encoding/json"
)

// Job is a command to run
type Job struct {
	Name         string
	Cmd          string
	WorkDir      string
	Args         []string
	MaxInstances int
	Timeout      JobDuration `json:",string"`
	RunsEvery    JobDuration `json:",string"`

	procMutex sync.Mutex
	Processes []*Process `json:",omitempty"`
}

func (j *Job) GetProcesses() []*Process {
	return j.Processes
}

func (j *Job) AddProcess(p *Process) error {
	j.procMutex.Lock()
	defer j.procMutex.Unlock()
	if len(j.Processes) > j.MaxInstances {
		log.Println("Max instances")
		return errors.New("Max instances")
	}
	j.Processes = append(j.Processes, p)
	return nil
}

func (j *Job) RemoveProcess(p *Process) error {
	j.procMutex.Lock()
	defer j.procMutex.Unlock()
	for k, v := range j.Processes {
		if v == p {
			j.Processes = append(j.Processes[:k], j.Processes[k+1:]...)
			break
		}
	}
	return nil
}

type JobDuration struct {
	time.Duration
}

func (jd *JobDuration) UnmarshalJSON(buf []byte) error {
	td, err := time.ParseDuration(strings.Trim(string(buf), `"`))
	if err != nil {
		return err
	}

	jd.Duration = td
	return nil
}

func (j *JobDuration) MarshalJSON() ([]byte, error) {
	return json.Marshal(j.Duration.String())
}