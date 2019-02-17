package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"os"
)



func SetWatcher(dirPath string,queueSize int){
	Queue = make(chan Job,queueSize)
	watchDog,err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Fprintf(os.Stderr,"%s",err.Error())
	}
	pipe := make(chan bool)
	go func() {
		for {
			select {
			// watch for events
			case event := <-watchDog.Events:
				e := Event{Op:uint32(event.Op),File:event.String()}
				j := Job{Event:e}
				Queue <- j
			case err := <-watchDog.Errors:
				fmt.Fprintf(os.Stderr,err.Error())
			}
		}
	}()
	watchDog.Add(dirPath)
	<-pipe
}