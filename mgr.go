package taskmgr

import (
	"fmt"
	"time"
)

const (
	DATE_LAYOUT = "2006-01-02 15:04:05"
)

type Task interface {
	Do() error
	String() string
	Begin()
	Benchmark() string
}

type Bench struct {
	ts1 int64
}

func (this *Bench) Begin() {
	this.ts1 = time.Now().Unix()
}

func (this *Bench) Benchmark() string {
	bench := time.Now().Unix() - this.ts1
	return humanTime(bench)
}

type TaskMgr struct {
	exit     chan bool
	monitor  chan<- string
	taskChan <-chan Task
	errChan  chan<- string
}

func NewMgr(monitor, errChan chan<- string, taskChan <-chan Task) *TaskMgr {
	return &TaskMgr{exit: make(chan bool, 1), monitor: monitor, taskChan: taskChan, errChan: errChan}
}

func (this *TaskMgr) Run() {
	for {
		select {
		case t := <-this.taskChan:
			select {
			case this.monitor <- fmt.Sprintf("task [%s] will handle!", t.String()):
			default:
			}
			go func() {
				t.Begin()
				err := t.Do()
				if err != nil {
					select {
					case this.errChan <- fmt.Sprintf("[%s]execute error! reason is %v", t.String(), err):
					case <-time.After(time.Second * 1):
					}
				}

				select {
				case this.monitor <- fmt.Sprintf("task [%s] benchmark %v!", t.String(), t.Benchmark()):
				default:
				}
			}()
		case <-this.exit:
			select {
			case this.monitor <- "task manager will be exit!":
			case <-time.After(time.Second * 1):
			}
			return
		}
	}
}

func (this *TaskMgr) Stop() {
	close(this.exit)
}

func humanTime(t int64) string {
	switch {
	case t < 60:
		return fmt.Sprintf("%d秒", t)
	case t > 60 && t < 3600:
		m := t / 60
		s := t % 60
		return fmt.Sprintf("%d分%d秒", m, s)
	case t > 3600:
		h := t / 3600
		m := (t - h*3600) / 60
		s := (t - h*3600) % 60
		return fmt.Sprintf("%d时%d分%d秒", h, m, s)
	}
	return ""
}
