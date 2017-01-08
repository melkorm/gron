package main

import (
	"flag"
	"log"
	"os"
	"sync"

	gron "github.com/melkorm/gron/lib"
)

func main() {
	f, err := os.OpenFile("logs", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	log.SetOutput(os.Stdout)
	path := flag.String("path", "", "Path, absolute or relative to  file")
	flag.Parse()

	c, err := os.Open(*path)

	if err != nil {
		log.Fatalln("%s", err)
		panic(err)
	}
	jsonConf := &gron.JSONConf{}
	tasks, err := jsonConf.Parse(c)

	if err != nil {
		log.Fatalf("Can't parse %s", err)
		panic(err)
	}
	var wg sync.WaitGroup

	spinner := gron.NewSpinner(tasks)

	spinner.Spin()
	go Serve(tasks)
	wg.Add(1)
	wg.Wait()
}
