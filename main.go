package main


func main() {
	sch := newScheduler()
	sch.addWorker(WorkerOne)
	sch.addWorker(WorkerTwo)
	sch.addWorker(WorkerThree)
	sch.run()
}
