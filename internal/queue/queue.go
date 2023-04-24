package queue

import (
	"sync"

	"LoLItemRecommender/internal/printer"
)

type JobTodo func() error
type Job struct {
	todo JobTodo
	next *Job
	prev *Job
}

func (j *Job) Do() error {
	return j.todo()
}

type Queue struct {
	m      *sync.RWMutex
	root   *Job
	last   *Job
	length int32
}

func NewQueue() *Queue {
	return &Queue{
		m: &sync.RWMutex{},
	}
}

const maxQueueLimitSize = 200

func (q *Queue) AddJob(j JobTodo) {
	q.m.Lock()
	defer q.m.Unlock()
	if q.length > maxQueueLimitSize {
		printer.Warn("Queue limit size reached, ignoring entry")
		return
	}
	nj := &Job{next: q.root, todo: j}
	if q.last == nil {
		q.last = nj
	}
	if q.root != nil {
		q.root.prev = nj
	}
	q.root = nj
	q.length++
}

func (q *Queue) Empty() bool {
	q.m.RLock()
	defer q.m.RUnlock()
	return q.length == 0
}
func (q *Queue) Size() int32 {
	q.m.RLock()
	defer q.m.RUnlock()
	return q.length
}

func (q *Queue) PopLastJob() *Job {
	q.m.Lock()
	defer q.m.Unlock()
	l := q.last
	if l == nil {
		return nil
	}
	q.last = q.last.prev
	q.length--
	return l
}
