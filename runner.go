package gron

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/twinj/uuid"
)

// Task is a task to run
type Task struct {
	Name             string
	Command          string
	Args             []string
	AllowedInstances int
	Instances        int
	TimeoutDuration  string
	RunsEvery        string
	Processes        []*Process
	duration         time.Duration
}

func (job *Task) Interval() time.Duration {
	if job.duration == 0 {
		duration, err := time.ParseDuration(job.RunsEvery)

		if err != nil {
			panic(err)
		}

		job.duration = duration
	}

	return job.duration
}

func (job *Task) Timeout() time.Duration {
	duration, err := time.ParseDuration(job.TimeoutDuration)
	if err != nil {
		panic(err)
	}

	return duration
}

func (t *Task) AddProcess(p *Process) error {
	if len(t.Processes) >= t.AllowedInstances {
		return errors.New("cant add new process")
	}

	t.Processes = append(t.Processes, p)

	return nil
}

func (t *Task) RemoveProcess(p *Process) error {
	log.Println("want to remove", p)
	for key, _ := range t.Processes {
		if t.Processes[key] == p {
			t.Processes = append(t.Processes[:key], t.Processes[key+1:]...)
			log.Println("removed", p)
			return nil
		}
	}

	return errors.New("process not found")
}

// FinishedJob type describing finished job
type Process struct {
	Start     time.Time
	End       time.Time
	ID        uuid.UUID
	Task      *Task
	State     string
	CmdReturn []byte
	Quit      chan int
}

func produce(t *Task, pipe chan<- *Process) {
	ticker := time.NewTicker(t.Interval())
	for range ticker.C {
		go func() {
			u := uuid.NewV4()
			p := &Process{ID: u, Task: t, State: "New"}
			err := t.AddProcess(p)
			if err == nil {
				log.Println(p.ID, p.State, p.Task.Name)
				pipe <- p
			}
		}()
	}
}

func consume(in chan *Process, out chan<- *Process) {
	for p := range in {
		go func(p *Process) {
			p.State = "To process"
			timeout := time.NewTicker(p.Task.Timeout())
			commandChannel := make(chan *Process, 1)
			p.Quit = make(chan int, 1)
			go runCommand(p, commandChannel)
			select {
			case <-commandChannel:
				p.End = time.Now()
				p.State = "Complete"
				p.Task.RemoveProcess(p)
				out <- p
				return
			case <-timeout.C:
				p.State = "Timeout"
				p.Task.RemoveProcess(p)
				out <- p
				return
			case <-p.Quit:
				p.State = "Killed"
				p.Task.RemoveProcess(p)
				out <- p
				return
			}
		}(p)
	}
}

func runCommand(p *Process, out chan *Process) {
	p.Start = time.Now()
	p.State = "Received"
	output, err := exec.Command(p.Task.Command, p.Task.Args...).Output()
	if err != nil {
		log.Fatalln("Command error", os.Stderr, err)
	}
	p.CmdReturn = output
	out <- p
}

func Run(tasks []*Task) chan *Process {
	out := make(chan *Process)
	in := make(chan *Process)
	for _, t := range tasks {
		go produce(t, in)
	}
	go consume(in, out)

	return out
}
