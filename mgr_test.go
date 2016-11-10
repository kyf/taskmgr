package taskmgr

import (
	"errors"
	"testing"
	"time"

	"github.com/kyf/util/log"
)

type DemoTask struct {
	Bench
}

func (this *DemoTask) Do() error {
	time.Sleep(time.Second * 5)
	return errors.New("demo testing error")
}

func (this *DemoTask) String() string {
	return "DemoTask testing string"
}

func TestMgr(t *testing.T) {
	monitor, errChan, taskChan := make(chan string, 1), make(chan string, 1), make(chan Task, 1)

	mgr := NewMgr(monitor, errChan, taskChan)
	go mgr.Run()

	exit := make(chan bool, 1)

	go func() {
		for {
			select {
			case it := <-monitor:
				log.Print(it)
			case it := <-errChan:
				log.Error(it)
			case <-exit:
				return
			}
		}
	}()

	for i := 0; i <= 10; i++ {
		taskChan <- &DemoTask{}
	}

	time.Sleep(time.Second * 10)
	mgr.Stop()
	time.Sleep(time.Second * 3)
	close(exit)
}
