package gron

import (
	"errors"
	"log"
	"strings"
	"sync"
	"time"
	"fmt"
)

// Job is a command to run
type Job struct {
	sync.Mutex
	Name         string
	Cmd          string
	Args         []string
	MaxInstances int
	Timeout      jobDuration
	RunsEvery    jobDuration
	Processes    []*Process
}

func (j *Job) GetProcesses() []*Process {
	return j.Processes
}

func (j *Job) AddProcess(p *Process) error {
	j.Lock()
	defer j.Unlock()
	if len(j.Processes) > j.MaxInstances {
		log.Println("Max instances")
		return errors.New("Max instances")
	}
	j.Processes = append(j.Processes, p)
	return nil
}

func (j *Job) RemoveProcess(p *Process) error {
	j.Lock()
	defer j.Unlock()
	for k, v := range j.Processes {
		if v == p {
			fmt.Println("pre removed", len(j.Processes))
			j.Processes = append(j.Processes[:k], j.Processes[k+1:]...)
			fmt.Println("post removed", len(j.Processes))
			break
		}
	}
	return nil
}

type jobDuration struct {
	time.Duration
}

func (jd *jobDuration) UnmarshalJSON(buf []byte) error {
	td, err := time.ParseDuration(strings.Trim(string(buf), `"`))
	if err != nil {
		return err
	}

	jd.Duration = td
	return nil
}
