package queue

import (
	"time"

	"LoLItemRecommender/internal/printer"
	"github.com/google/uuid"
)

type Pool struct {
	q                  *Queue
	cap                int
	chanStopInProgress chan bool
	chanErr            chan error
	chanWorkerDone     chan *Worker
	workers            map[string]*Worker
}

func (p *Pool) GetErrorsChan() <-chan error {
	return p.chanErr
}

func (p *Pool) Close() {
	p.chanStopInProgress <- true
	for id, w := range p.workers {
		printer.Debug("Worker %s stopped", id)
		w.Stop()
	}
	close(p.chanStopInProgress)
	close(p.chanErr)
	close(p.chanWorkerDone)
}

func (p *Pool) Dispatch(j JobTodo) {
	w := p.GetFirstAsleepWorker()
	if w != nil {
		w.ChangeState(StatePending)
		printer.Debug("[Dispatch] Attributing to %p", w)
		w.Attribute(j)
	} else {
		printer.Debug("[Dispatch] No worker available for this job, adding to the queue (%p)", j)
		p.q.AddJob(j)
	}
}

func (p *Pool) GetFirstAsleepWorker() *Worker {
	for _, w := range p.workers {
		if w.GetState() == StateAsleep {
			return w
		}
	}
	return nil
}

func (p *Pool) dequeueJob(initEnded chan<- bool) {
	initEnded <- true
	for {
		select {
		case w := <-p.chanWorkerDone:
			j := p.q.PopLastJob()
			if j != nil {
				w.ChangeState(StatePending)
				w.Attribute(j.todo)
			}
		case <-p.chanStopInProgress:
			return
		}
	}
}

func (p *Pool) allWorkerAreDone() bool {
	for _, w := range p.workers {
		if w.GetState() != StateAsleep {
			return false
		}
	}
	return true
}

func (p *Pool) WaitJobsToComplete() {
	for !p.q.Empty() || !p.allWorkerAreDone() {
		time.Sleep(1 * time.Second)
	}
}

func NewPool(cap int) *Pool {
	printer.Info("Creating pool with a capacity of {-F_BLUE,BOLD}%d {-RESET}worker(s)", cap)
	p := &Pool{
		q:              NewQueue(),
		cap:            cap,
		workers:        make(map[string]*Worker, cap),
		chanErr:        make(chan error),
		chanWorkerDone: make(chan *Worker),
	}
	rdyQueueJob := make(chan bool)
	for i := 0; i < cap; i++ {
		id := uuid.NewString()
		p.workers[id] = NewWorker(p.chanErr, p.chanWorkerDone)
		printer.Info("Worker {-F_GREEN}%s{-RESET} (%p) ready", id, p.workers[id])
	}
	go p.dequeueJob(rdyQueueJob)
	<-rdyQueueJob
	return p
}

func CalculatePoolCap[T any](jobs []T) int {
	c := len(jobs) / 3
	if c == 0 {
		return 1
	}
	return c
}
