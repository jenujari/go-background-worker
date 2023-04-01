package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type WorkerFn func(*JobScheduler)

type JobScheduler struct {
	interuppted bool
	interrupt   chan os.Signal

	ctx    context.Context
	cancel context.CancelFunc

	wg   *sync.WaitGroup
	done chan bool
	fatalErrorChan chan error

	workers []WorkerFn
}

func newScheduler() *JobScheduler {
	job := new(JobScheduler)
	job.ctx, job.cancel = context.WithCancel(context.Background())
	job.wg = new(sync.WaitGroup)
	job.done = make(chan bool)
	job.fatalErrorChan = make(chan error)
	job.interrupt = make(chan os.Signal)
	job.interuppted = false

	return job
}

func (j *JobScheduler) addWorker(fn WorkerFn) {
	j.workers = append(j.workers, fn)
}

func (j *JobScheduler) run() {
	if len(j.workers) == 0 {
		log.Fatalln("No workers added before invoking run.")
	}

	for _ , wrkr := range j.workers {
		j.wg.Add(1)
		go wrkr(j)
	}

	j.waitForFinish()
}

func (j *JobScheduler) waitForFinish() {
	go handleInterrupt(j)
	go waitGroupDone(j)
	go watchError(j)
	gracefullExit(j)
}

func handleInterrupt(j *JobScheduler) {
	signal.Notify(j.interrupt, syscall.SIGINT, syscall.SIGTERM)
	for range j.interrupt {
		if j.interuppted {
			log.Println("\nInterrupt signal already captured working on closing the process.")
			continue
		}
		j.interuppted = true
		j.cancel()
		fmt.Println("\nInteruppt signal captured.")
	}
}

func watchError(j *JobScheduler) {
	err := <- j.fatalErrorChan
	log.Println("Fatal error captured :: ", err)
	j.cancel()
}

func waitGroupDone(j *JobScheduler) {
	j.wg.Wait()
	j.done <- true
}

func gracefullExit(j *JobScheduler) {
	<-j.done
	fmt.Println("Gracefull exit")
	os.Exit(0)
}
