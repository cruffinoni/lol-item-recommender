package queue

import (
	"sync"

	"LoLItemRecommender/internal/printer"
)

type WorkerState uint8

const (
	StateAsleep = iota
	StatePending
	StateWorking
)

type Worker struct {
	r              *sync.RWMutex
	state          WorkerState
	errorChan      chan<- error
	workDone       chan<- *Worker
	scheduler      chan JobTodo
	terminate      chan bool
	terminateState bool
}

func (w *Worker) ListenJobs() {
	for {
		select {
		case <-w.terminate:
			return
		case j := <-w.scheduler:
			if w.terminateState {
				return
			}
			printer.Debug("[%p] Job received: %p", w, j)
			printer.Debug("[%p] Changing state to working", w)
			w.ChangeState(StateWorking)
			if err := j(); err != nil {
				printer.Debug("[%p] Sending err '%s'", w, err.Error())
				w.errorChan <- err
			}
			w.ChangeState(StateAsleep)
			w.workDone <- w
			printer.Debug("[%p] Changing state to asleep", w)
			if w.terminateState {
				return
			}
		}
	}
}

func (w *Worker) Stop() {
	w.terminate <- true
	w.terminateState = true
	close(w.errorChan)
	close(w.workDone)
	close(w.scheduler)
	close(w.terminate)
}

func (w *Worker) ChangeState(s WorkerState) {
	w.r.Lock()
	defer w.r.Unlock()
	w.state = s
}

func (w *Worker) GetState() WorkerState {
	w.r.RLock()
	defer w.r.RUnlock()
	return w.state
}

func (w *Worker) Attribute(j JobTodo) {
	w.scheduler <- j
}

func NewWorker(errorChan chan<- error, workDone chan<- *Worker) *Worker {
	w := &Worker{
		state:     StateAsleep,
		errorChan: errorChan,
		scheduler: make(chan JobTodo),
		terminate: make(chan bool),
		workDone:  workDone,
		r:         &sync.RWMutex{},
	}
	go w.ListenJobs()
	return w
}
