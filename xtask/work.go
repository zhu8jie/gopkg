package xtask

import "math/rand"

type Worker struct {
	taskChanel []chan Task
	workCnt    int
}

func NewWorker(workCnt, chanelCnt int) *Worker {
	wk := new(Worker)
	for i := 0; i < workCnt; i++ {
		wk.taskChanel = append(wk.taskChanel, make(chan Task, chanelCnt))
	}
	wk.workCnt = workCnt
	return wk
}

func (wk *Worker) Start() {
	for i := 0; i < wk.workCnt; i++ {
		go wk.workerThread(i)
	}
}

func (wk *Worker) workerThread(index int) {
	for {
		task := <-wk.taskChanel[index]
		task.DoSomeThingFirst()
		task.DoSomeThingSecond()
		task.DoSomeThingThird()
	}
}

func (wk *Worker) AddTask(tsk Task) {
	wk.taskChanel[rand.Intn(wk.workCnt)] <- tsk
}
