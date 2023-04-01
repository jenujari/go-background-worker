package main

import (
	"errors"
	"fmt"
	"time"
)

func WorkerOne(j *JobScheduler) {
	for {
		select {
		case <-j.ctx.Done():
			{
				fmt.Println("context closed exiting process 1.")
				j.wg.Done()
				return
			}
		default:
			{
				time.Sleep(time.Second * 2)
				fmt.Println("background worker one")
			}
		}
	}
}

func WorkerTwo(j *JobScheduler) {
	for {
		select {
		case <-j.ctx.Done():
			{
				fmt.Println("context closed exiting process 2.")
				j.wg.Done()
				return
			}
		default:
			{
				time.Sleep(time.Second* 4)
				fmt.Println("background worker two")
			}
		}
	}
}

func WorkerThree(j *JobScheduler) {
	var count int = 0
	for {
		select {
		case <-j.ctx.Done():
			{
				fmt.Println("context closed exiting process 3.")
				j.wg.Done()
				return
			}
		default:
			{
				time.Sleep(time.Second* 6)
				fmt.Println("background worker three")
				count++
				if count > 3 {
						j.fatalErrorChan <- errors.New("manual error occured from third worker") 
						j.wg.Done()
						return
				}
			}
		}
	}
}
