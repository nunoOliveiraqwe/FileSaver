package main

import (
	"fmt"
	"github.com/mediocregopher/radix.v2/pool"
	"github.com/mediocregopher/radix.v2/redis"
	"io/ioutil"
	"os"
	"strings"
)

type Job struct {
	Event Event
}
var Queue chan Job

type Event struct {
	Op uint32
	File string
}

type Worker struct {
	WorkerPool chan chan Job
	JobChannel chan Job
	Stop chan bool
}

var DB *pool.Pool

func NewWorker(workerPool chan chan Job) Worker{
	return Worker{
		WorkerPool:workerPool,
		JobChannel:make(chan Job),
		Stop: make(chan bool),
	}
}

func (w Worker) Start(){
	go func(){
		for{
			w.WorkerPool <- w.JobChannel
			select {
			case job := <- w.JobChannel:
				w.parseEvent(job.Event)
			case <-w.Stop:
				return

			}
		}
	}()
}

func (w Worker) Quit(){
	w.Stop <- true
}

func (w Worker) parseEvent(evt Event){
	if evt.Op == 1 {
		filePathSlice := strings.Split(evt.File,"\"")
		if len(filePathSlice)>1{
			filePath := filePathSlice[1]
			contents, err := w.readFile(filePath)
			if err != nil {
				fmt.Fprintf(os.Stderr,err.Error())
			} else{
				fileNameSlice := strings.Split(filePath,string(os.PathSeparator))
				if len(fileNameSlice) > 0 {
					filename := fileNameSlice[len(fileNameSlice)-1]
					w.parseToRedis(contents,filename)
				} else{
					fmt.Fprintf(os.Stderr,"Filename impossible to determine!! Something is wrong")
				}
			}
		}
	}
}

func (w Worker) readFile(path string) (string,error){
	file, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr,err.Error())
		return "",err
	} else{
		return string(file),nil
	}
}


func (w Worker) parseToRedis(content string, filename string){
	conn, err := DB.Get()
	if err != nil {
		fmt.Fprint(os.Stderr,"Impossible to insert file %s\nError %s\n",filename,err.Error())
	} else{
		defer DB.Put(conn)
		postSlice := strings.Split(filename,";")
		var reply *redis.Resp
		if len(postSlice)==2 {
			var check = false
			if strings.ToLower(postSlice[0]) == "insert" || strings.ToLower(postSlice[0]) == "update" {
				time := GetNtpTime()
				reply = conn.Cmd("HMSET", "blogPost:"+postSlice[1], "title", postSlice[1], "content", content, "timestamp", time)
				check=true
			} else if strings.ToLower(postSlice[0]) == "delete"{
				reply = conn.Cmd("HDEL","blogPost:"+postSlice[1])
				check = true
			} else{
				fmt.Fprintf(os.Stderr,"Invalid operation command %s\n",postSlice[0])
			}
			if check==true && reply.Err != nil {
				fmt.Fprintf(os.Stderr,err.Error())
			} else if check==true{
				fmt.Fprintf(os.Stdout,"Sucefully insert %s\n",filename)
			}
		} else{
			fmt.Fprint(os.Stderr,"Invalid file %s\n",filename)
		}

	}
}