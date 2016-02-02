package main

import (
	"flag"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/melkorm/gron/config"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
}

// FinishedJob type describing finished job
type FinishedJob struct {
	Job    config.Job
	Output string
}

func runJob(job config.Job, finishedJobs chan FinishedJob) {
	ticker := time.NewTicker(job.Interval())
	go func() {
		for range ticker.C {
			if job.CanRun() {
				job.IncrementInstances()
				go func() {
					ouput, err := exec.Command(job.Command, job.Args...).Output()
					if err != nil {
						log.Fatalln("Command error", os.Stderr, err)
					}

					finishedJobs <- FinishedJob{job, string(ouput) + time.Now().String()}
					job.DecrementInstances()
				}()
			} else {
				log.Debug("Job " + job.Name + " exceeded runners " + string(job.Instances))
			}
		}
	}()
}

func main() {
	path := flag.String("path", "", "Path, absolute or relative to config file")
	flag.Parse()

	data, err := ioutil.ReadFile(*path)

	if err != nil {
		log.Fatalln("%s", err)
		panic(err)
	}
	test, err := config.Parse(data)

	if err != nil {
		log.Fatalln("Cant parse config %s", err)
		panic(err)
	}

	finishedJobs := make(chan FinishedJob)

	for _, job := range test.Jobs {
		go runJob(job, finishedJobs)
	}

	for finishedJob := range finishedJobs {
		log.Debug("Finished", finishedJob.Job.Name, "with output", finishedJob.Output)
	}
}
