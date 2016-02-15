package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/melkorm/gron/config"
	"github.com/melkorm/gron/rest"
	"github.com/melkorm/gron/runner"
)

func main() {
	f, err := os.OpenFile("logs", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	log.SetOutput(os.Stdout)
	// log := log.New(os.Stdout, "new", 1)

	path := flag.String("path", "", "Path, absolute or relative to  file")
	flag.Parse()

	data, err := ioutil.ReadFile(*path)

	if err != nil {
		log.Fatalln("%s", err)
		panic(err)
	}
	tasks, err := config.Parse(data)

	if err != nil {
		log.Fatalf("Can't parse %s", err)
		panic(err)
	}

	runner.Run(tasks)
	rest.Serve(tasks)
}
