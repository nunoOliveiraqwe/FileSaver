package main

import (
	"flag"
	"fmt"
	"github.com/mediocregopher/radix.v2/pool"
	"os"
	"strings"
)

func main(){
	var directoryPath string
	var redisServerAddr string
	var redisConnectionType string
	flag.StringVar(&directoryPath,"directory","","Path to the directory to monitor")
	flag.StringVar(&redisServerAddr,"redisServerAddr","","Redis Server Address")
	flag.StringVar(&redisConnectionType,"redisConnectionType","","Redis Connection Type: TCP OR UDP")
	var maxWorkers = flag.Int("maxWorkers",10,"Maximum Number of Workers")
	var required = []string{"directory","redisServerAddr","redisConnectionType"}
	flag.Parse()
	seen := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) { seen[f.Name] = true })
	for _, req := range required {
		if !seen[req] {
			fmt.Fprintf(os.Stderr, "missing required -%s argument/flag\n", req)
			os.Exit(2)
		}
	}
	if strings.ToLower(redisConnectionType) != "tcp" && strings.ToLower(redisConnectionType) != "udp"{
		fmt.Fprintf(os.Stderr,"invalid argument for redisConnectionType")
		os.Exit(2)
	}
	initRedisConnection(redisServerAddr,strings.ToLower(redisConnectionType),*maxWorkers)
	Start(directoryPath,*maxWorkers)
}


func initRedisConnection(redisServerAddr string,redisConnectionType string,maxWorkers int){
	var err error
	DB,err = pool.New(redisConnectionType,redisServerAddr,maxWorkers)
	if err != nil{
		fmt.Fprint(os.Stderr,err.Error())
		os.Exit(2)
	}
}

func Start(directoryPath string,maxWorkers int){
	dispatcher := NewDispatcher(maxWorkers)
	dispatcher.Start(maxWorkers)
	SetWatcher(directoryPath,maxWorkers)
}