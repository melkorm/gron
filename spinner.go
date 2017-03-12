package gron

import (
	"context"
	"fmt"
	"os/exec"
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
func (s *Spinner) Spin() {
	for _, j := range s.jobs {
		go func(j *Job) {
			// to not spam all the jobs at once
			// sleep second between
			time.Sleep(1 * time.Second)
			ticker := time.NewTicker(j.RunsEvery.Duration)
			instances := make(chan int, j.MaxInstances)
			for range ticker.C {
				instances <- 1
				timeoutCtx, _ := context.WithTimeout(context.Background(), j.Timeout.Duration)
				cmd := exec.CommandContext(timeoutCtx, j.Cmd, j.Args...)
				cmd.Dir = j.WorkDir
				p := &Process{Id: uuid.NewV4().String(), Cmd: cmd}
				// err should never happen
				err := j.AddProcess(p)
				if err == nil {
					go func(p *Process, j *Job, instances <-chan int, timeoutCtx context.Context) {
						p.run()
						if timeoutCtx.Err() != nil {
							fmt.Println("Context error:", timeoutCtx.Err())
						}
						j.RemoveProcess(p)
						<-instances
					}(p, j, instances, timeoutCtx)
				} else {
					fmt.Println(err)
				}
			}
		}(j)
	}
}
