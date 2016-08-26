package gron

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/twinj/uuid"
)

type Spinner struct {
	wg   *sync.WaitGroup
	jobs []*Job
	log  *log.Logger
}

func NewSpinner(jobs []*Job, wg *sync.WaitGroup, log *log.Logger) *Spinner {
	return &Spinner{jobs: jobs, wg: wg, log: log}
}

// Spin all jobs
func (s *Spinner) Spin() error {
	for _, j := range s.jobs {
		go func(j *Job) {
			ticker := time.NewTicker(j.RunsEvery.Duration)
			for range ticker.C {
				cmd := exec.Command(j.Cmd, j.Args...)
				p := &Process{id: uuid.NewV4().String(), cmd: cmd}
				j.proccesses = append(j.proccesses, p)
				go p.run()
			}
		}(j)
	}
	return nil
}

// Process type describes single task's process
type Process struct {
	id     string
	cmd    *exec.Cmd
	err    error
	output []byte
	state  string
	start  time.Time
	finish time.Time
}

func (p *Process) run() {
	p.start = time.Now()
	log.Println("Started", p)
	p.output, p.err = p.cmd.CombinedOutput()
	p.finish = time.Now()
	log.Println("Finished", p)
}

func (p *Process) String() string {
	return fmt.Sprintf("ID: %s; Command: %s; Started: %s; Finished %s", p.id, p.cmd.Args, p.start, p.finish)
}

// Job is a command to run
type Job struct {
	Name         string
	Cmd          string
	Args         []string
	MaxInstances int
	Timeout      jobDuration
	RunsEvery    jobDuration
	proccesses   []*Process
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
