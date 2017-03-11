package gron

import (
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/twinj/uuid"
	"os"
)

type Spinner struct {
	jobs []*Job
}

func NewSpinner(jobs []*Job) *Spinner {
	return &Spinner{jobs: jobs}
}

// Spin all jobs
func (s *Spinner) Spin() error {
	for _, j := range s.jobs {
		go func(j *Job) {
			ticker := time.NewTicker(j.RunsEvery.Duration)
			for range ticker.C {
				cmd := exec.Command(j.Cmd, j.Args...)
				p := &Process{Id: uuid.NewV4().String(), cmd: cmd}
				err := j.AddProcess(p)
				if err == nil {
					go func(p *Process, j *Job) {
						p.run()
						j.RemoveProcess(p)
					}(p, j)
				}
			}
		}(j)
	}
	return nil
}

// Process type describes single task's process
type Process struct {
	Id       string
	cmd      *exec.Cmd
	err      error
	Output   []byte
	State    string
	Start    time.Time
	Finished time.Time
}

func (p *Process) run() {
	p.Start = time.Now()
	log.Println("Started", p)
	p.Output, p.err = p.cmd.CombinedOutput()
	p.Finished = time.Now()
	log.Println("Finished", p)
}

func (p *Process) String() string {
	return fmt.Sprintf("ID: %s; Command: %s; Started: %s; Finished %s; Err: %v; Output: %v", p.Id, p.cmd.Args, p.Start, p.Finished, p.err, string(p.Output))
}

func (p Process) GetID() string {
	return p.Id
}

func (p Process) Signal(sig os.Signal) error {
	return p.cmd.Process.Signal(sig)
}