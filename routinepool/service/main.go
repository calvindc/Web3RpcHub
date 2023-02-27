package main

import (
	"fmt"
	"runtime"
	"time"
)

//仅仅是个框架例子--------------------------------------------

type MessageProc struct {
	MsgHash int
}

func (m *MessageProc) Do() {
	fmt.Println("MegHash Num:", m.MsgHash)
	//do something
	time.Sleep(100 * time.Millisecond)
}

func main() {
	//开启两万个县城
	routimeNum := 100 * 100 * 2
	p := NewWorkerpool(routimeNum)
	p.Run()

	go func() {
		for i := 1; i < 100*100*100*10; i++ {
			mh := &MessageProc{
				MsgHash: i,
			}

			p.jobQu <- mh
		}
	}()

	for {
		fmt.Println(fmt.Sprintf("runtime.NumGoroutine= %d ,NumCPU = %d ,NumCgoCall = %d", runtime.NumGoroutine(), runtime.NumCPU(), runtime.NumCgoCall()))
		time.Sleep(1 * time.Second)
	}

}

//-----------------job-----------------
type Job interface {
	Do()
}
type JobQu chan Job

//-----------------worker---------------
type Worker struct {
	jobChan JobQu //每个work对象具备JobQu属性
}

func NewWoker() Worker {
	return Worker{
		jobChan: make(chan Job),
	}
}

func (w *Worker) Run(wq chan JobQu) {
	go func() {
		for {
			wq <- w.jobChan
			select {
			case job := <-w.jobChan: //-------------【读】，读写是统一数据，worker不会死锁
				jo := job.(Job) //对各种任务进行断言
				fmt.Println("==>执行任务的worker地址：", w.jobChan)
				jo.Do()
			}
		}
	}()
}

//-----------------worker Pool-----------
type Workerpool struct {
	WorkerLength int        //容量
	jobQu        JobQu      //job队列 接收外部数据
	WorkQu       chan JobQu //worker队列 处理work任务的
}

func NewWorkerpool(worklen int) *Workerpool {
	return &Workerpool{
		WorkerLength: worklen,
		jobQu:        make(JobQu),
		WorkQu:       make(chan JobQu, worklen),
	}
}

func (wp *Workerpool) Run() {
	fmt.Println("init hub worker")
	for i := 0; i < wp.WorkerLength; i++ {
		worker := NewWoker()
		worker.Run(wp.WorkQu)
	}

	go func() {
		for {
			select {
			case job := <-wp.jobQu: //协程池中呆处理的work,来自于jobQu
				worker := <-wp.WorkQu
				worker <- job //空闲的线程执行-------------【写】
				fmt.Println("临时worker地址：", worker)
				fmt.Println(fmt.Sprintf("job类型%T, worker类型%T", job, worker))
			}
		}
	}()
}
