// Package main  provides ...
package main

import (
	"chat"
	"flag"
	"log"
	"os"
	"runtime/pprof"
	"time"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {

	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)

		go func() {
			time.Sleep(2 * time.Minute)
			pprof.StopCPUProfile()
			log.Printf("Stoped profile")
			os.Exit(0)
		}()
	}

	server := &chat.ChatServer{":12345", make(map[string]*chat.Room), make(map[string]*chat.Client)}
	server.ListenAndServe()
}
