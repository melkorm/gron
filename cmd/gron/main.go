package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/melkorm/gron"
)

func main() {
	log.SetOutput(os.Stdout)
	path := flag.String("path", "", "Path, absolute or relative to  file")
	flag.Parse()

	c, err := os.Open(*path)

	if err != nil {
		log.Fatalln("%s", err)
		panic(err)
	}
	defer c.Close()

	jsonConf := &gron.JSONConf{}
	tasks, err := jsonConf.Parse(c)

	if err != nil {
		log.Fatalf("Can't parse %s", err)
		panic(err)
	}
	var wg sync.WaitGroup
	spinner := gron.NewSpinner(tasks)
	go gron.Serve(tasks)
	go spinner.Spin()
	wg.Add(1)
	go signalHandler(wg)
	wg.Wait()
}

func signalHandler(group sync.WaitGroup) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until a signal is received.
	s := <-c
	fmt.Println("Got signal:", s)
	if s == os.Interrupt {
		fmt.Println("Attempting gracefull shut down ...")
		defer group.Done()
		os.Exit(0)
	}
	signal.Stop(c)
}
