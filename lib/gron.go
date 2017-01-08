package gron

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/twinj/uuid"
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
				SigChan := make(chan os.Signal)
				p := &Process{id: uuid.NewV4().String(), cmd: cmd, SigChan: SigChan}
				err := j.AddProcess(p)
				if err == nil {
					go p.run()
				}
			}
		}(j)
	}
	return nil
}

// Process type describes single task's process
type Process struct {
	id      string
	cmd     *exec.Cmd
	err     error
	output  []byte
	state   string
	start   time.Time
	finish  time.Time
	SigChan chan os.Signal
}

func (p *Process) run() {
	go func(p *Process) {
		p.start = time.Now()
		log.Println("Started", p)
		p.output, p.err = p.cmd.CombinedOutput()
		p.finish = time.Now()
		log.Println("Finished", p)
	}(p)
	for sig := range p.SigChan {
		log.Println("Received signal", sig)
		if p.cmd.Process != nil {
			p.cmd.Process.Signal(sig)
		}
	}
}

func (p *Process) String() string {
	return fmt.Sprintf("ID: %s; Command: %s; Started: %s; Finished %s; Err: %s; Output: %s", p.id, p.cmd.Args, p.start, p.finish, p.err, string(p.output))
}

func (p Process) GetID() string {
	return p.id
}

// Job is a command to run
type Job struct {
	sync.Mutex
	Name         string
	Cmd          string
	Args         []string
	MaxInstances int
	Timeout      jobDuration
	RunsEvery    jobDuration
	proccesses   []*Process
}

func (j *Job) GetProcesses() []*Process {
	return j.proccesses
}

func (j *Job) AddProcess(p *Process) error {
	j.Lock()
	defer j.Unlock()
	if len(j.proccesses) > j.MaxInstances {
		log.Printf("Max instances")
		return errors.New("Max instances")
	}
	j.proccesses = append(j.proccesses, p)
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
