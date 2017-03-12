package gron

import (
	"log"
	"os"
	"os/exec"
	"time"
	"encoding/json"
)

// Process type describes single task's process
type Process struct {
	Id       string
	Cmd      *exec.Cmd
	Err      error
	Output   string
	Started  time.Time
	Finished time.Time
}

func (p *Process) run() {
	p.Started = time.Now()
	m, _ := json.Marshal(p)
	log.Println("Started", string(m))
	out, err := p.Cmd.CombinedOutput()
	p.Output = string(out)
	p.Err = err
	p.Finished = time.Now()
	m, _ = json.Marshal(p)
	log.Println("Finished", string(m))
}

func (p Process) Signal(sig os.Signal) error {
	return p.Cmd.Process.Signal(sig)
}