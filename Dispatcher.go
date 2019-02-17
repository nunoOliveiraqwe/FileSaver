package main

type Dispatcher struct {
	WorkerPool chan chan Job
}

func NewDispatcher(maxWorkers int) *Dispatcher{
	pool := make(chan chan Job,maxWorkers)
	return &Dispatcher{WorkerPool:pool}
}

func (d *Dispatcher) Start(maxWorkers int){
	for i:=0; i< maxWorkers; i++{
		worker := NewWorker(d.WorkerPool)
		worker.Start()
	}
	go d.dispatch()
}

func (d *Dispatcher) dispatch(){
	for{
		select {
		case job := <- Queue:
			go func(job Job) {
				channel := <- d.WorkerPool
				channel <- job
			}(job)
		}
	}
}