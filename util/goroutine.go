package util

import "sync"

// GoFunc runs a goroutine under WaitGroup
func GoFunc(routinesGroup *sync.WaitGroup, f func()) {
	routinesGroup.Add(1)
	go func() {
		defer routinesGroup.Done()
		f()
	}()
}

type WaitGroupWrapper struct {
	sync.WaitGroup
}

func (w *WaitGroupWrapper) Wrap(cb func()) {
	w.Add(1)
	go func() {
		defer w.Done()
		cb()
	}()
}
