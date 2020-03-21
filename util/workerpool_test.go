package util

import (
	"testing"
	"time"

	"github.com/tianhongw/misc-go/log"
	"github.com/tianhongw/misc-go/util/assert"
)

func TestWorkerPool(t *testing.T) {
	workPool := NewWorkerPool()
	if err := workPool.Run(f1); err != nil {
		t.Fail()
	}

	if err := workPool.Run(f2); err != nil {
		t.Fail()
	}

	workPool.Close()

	assert.Equal(t, ErrWorkerPoolClosed, workPool.Run(f1))
}

func f1() {
	t := time.NewTimer(1 * time.Second)
loop:
	for {
		select {
		case <-t.C:
			log.Instance().Info("workerPool: f1() done")
			break loop
		default:
		}
	}
}

func f2() {
	t := time.NewTimer(2 * time.Second)
loop:
	for {
		select {
		case <-t.C:
			log.Instance().Info("workerPool: f2() done")
			break loop
		default:
		}
	}
}
